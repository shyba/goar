package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/liteseed/goar/tx"
	"github.com/liteseed/goar/types"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
)

const (
	DEFAULT_CHUNK_CONCURRENT_NUM = 50
	MAX_CHUNKS_IN_BODY           = 1
)

type SerializedUploader struct {
	chunkIndex         int
	txPosted           bool
	transaction        *types.Transaction
	lastRequestTimeEnd int64
	lastResponseStatus int
	lastResponseError  string
}

type TransactionUploader struct {
	Client             *Client `json:"-"`
	ChunkIndex         int
	TxPosted           bool
	Transaction        *types.Transaction
	Data               []byte
	LastRequestTimeEnd int64
	TotalErrors        int // Not serialized.
	LastResponseStatus int
	LastResponseError  string
}

func newUploader(t *types.Transaction, client *Client) (*TransactionUploader, error) {
	if t.ID == "" {
		return nil, errors.New("Transaction is not signed.")
	}
	if t.ChunkData == nil {
		log.Println("Transaction chunks not prepared")
	}
	// Make a copy of Transaction, zeroing the Data so we can serialize.
	u := &TransactionUploader{
		Client: client,
	}

	u.Transaction = &types.Transaction{
		Format:    t.Format,
		ID:        t.ID,
		LastTx:    t.LastTx,
		Owner:     t.Owner,
		Tags:      t.Tags,
		Target:    t.Target,
		Quantity:  t.Quantity,
		Data:      t.Data,
		DataSize:  t.DataSize,
		DataRoot:  t.DataRoot,
		Reward:    t.Reward,
		Signature: t.Signature,
		ChunkData: t.ChunkData,
	}
	return u, nil
}

// CreateUploader
// @param upload: Transaction | SerializedUploader | string,
// @param Data the Data of the Transaction. Required when resuming an upload.
func CreateUploader(api *Client, upload interface{}, data []byte) (*TransactionUploader, error) {
	var (
		uploader *TransactionUploader
		err      error
	)

	if tt, ok := upload.(*types.Transaction); ok {
		uploader, err = newUploader(tt, api)
		if err != nil {
			return nil, err
		}
		return uploader, nil
	}

	if id, ok := upload.(string); ok {
		// upload 返回为 SerializedUploader 类型
		upload, err = (&TransactionUploader{Client: api}).FromTransactionId(id)
		if err != nil {
			log.Println("(&TransactionUploader{Client: api}).FromTransactionId(id)", "err", err)
			return nil, err
		}
	} else {
		// 最后 upload 为 SerializedUploader type
		newUpload, ok := upload.(*SerializedUploader)
		if !ok {
			panic("upload params error")
		}
		upload = newUpload
	}

	uploader, err = (&TransactionUploader{Client: api}).FromSerialized(upload.(*SerializedUploader), data)
	return uploader, err
}

func (tt *TransactionUploader) Once() (err error) {
	for !tt.IsComplete() {
		if err = tt.UploadChunk(); err != nil {
			return
		}

		if tt.LastResponseStatus != 200 {
			return errors.New(tt.LastResponseError)
		}
	}

	return
}

func (tt *TransactionUploader) IsComplete() bool {
	tChunks := tt.Transaction.ChunkData
	if tChunks == nil {
		return false
	} else {
		return tt.TxPosted && (tt.ChunkIndex == len(tChunks.Chunks)) || tt.TxPosted && len(tChunks.Chunks) == 0
	}
}

func (tt *TransactionUploader) TotalChunks() int {
	if tt.Transaction.ChunkData == nil {
		return 0
	} else {
		return len(tt.Transaction.ChunkData.Chunks)
	}
}

func (tt *TransactionUploader) UploadedChunks() int {
	return tt.ChunkIndex
}

func (tt *TransactionUploader) PctComplete() float64 {
	val := decimal.NewFromInt(int64(tt.UploadedChunks())).Div(decimal.NewFromInt(int64(tt.TotalChunks())))
	fval, _ := val.Float64()
	return math.Trunc(fval * 100)
}

func (tt *TransactionUploader) ConcurrentOnce(ctx context.Context, concurrentNum int) error {
	// post tx info
	if err := tt.postTransaction(); err != nil {
		return err
	}

	if tt.IsComplete() {
		return nil
	}

	var wg sync.WaitGroup
	if concurrentNum <= 0 {
		concurrentNum = DEFAULT_CHUNK_CONCURRENT_NUM
	}
	p, _ := ants.NewPoolWithFunc(concurrentNum, func(i interface{}) {
		defer wg.Done()
		// process submit chunk
		idx := i.(int)

		select {
		case <-ctx.Done():
			log.Println("ctx.done", "chunkIdx", idx)
			return
		default:
		}
		var chunk *types.ChunkData
		var err error

		if err != nil {
			log.Println("GetChunk error", "err", err, "idx", idx)
			return
		}
		body, statusCode, err := tt.Client.SubmitChunkData(chunk) // always body is errMsg
		if statusCode == 200 {
			return
		}

		log.Println("concurrent submitChunk failed", "chunkIdx", idx, "statusCode", statusCode, "gatewayErr", body, "httpErr", err)
		// try again
		retryCount := 0
		for {
			select {
			case <-ctx.Done():
				log.Println("ctx.done", "chunkIdx", idx)
				return
			default:
			}

			retryCount++
			if statusCode == 429 {
				time.Sleep(1 * time.Second)
			} else {
				time.Sleep(200 * time.Millisecond)
			}

			body, statusCode, err = tt.Client.SubmitChunkData(chunk)
			if statusCode == 200 {
				return
			}
			log.Println("retry submitChunk failed", "retryCount", retryCount, "chunkIdx", idx, "statusCode", statusCode, "gatewayErr", body, "httpErr", err)
		}
	})

	defer p.Release()
	for i := 0; i < len(tt.Transaction.ChunkData.Chunks); i++ {
		wg.Add(1)
		if err := p.Invoke(i); err != nil {
			log.Println("p.Invoke(i)", "err", err, "i", i)
			return err
		}
	}

	wg.Wait()
	return nil
}

/**
 * Uploads the next part of the Transaction.
 * On the first call this posts the Transaction
 * itself and on any subsequent calls uploads the
 * next chunk until it completes.
 */
func (tt *TransactionUploader) UploadChunk() error {
	return nil
}

/**
 * Reconstructs an upload from its serialized state and data.
 * Checks if data matches the expected data_root.
 *
 * @param serialized
 * @param data
 */
func (tt *TransactionUploader) FromSerialized(serialized *SerializedUploader, data []byte) (*TransactionUploader, error) {
	if serialized == nil {
		return nil, errors.New("Serialized object does not match expected format.")
	}

	// Everything looks ok, reconstruct the TransactionUpload,
	// prepare the chunks again and verify the data_root matches
	upload, err := newUploader(serialized.transaction, tt.Client)
	if err != nil {
		return nil, err
	}
	// Copy the serialized upload information, and Data passed in.
	upload.ChunkIndex = serialized.chunkIndex
	upload.LastRequestTimeEnd = serialized.lastRequestTimeEnd
	upload.LastResponseError = serialized.lastResponseError
	upload.LastResponseStatus = serialized.lastResponseStatus
	upload.TxPosted = serialized.txPosted
	upload.Data = data

	err = tx.PrepareChunks(upload.Transaction, data)
	if err != nil {
		return nil, err
	}

	if upload.Transaction.DataRoot != serialized.transaction.DataRoot {
		return nil, errors.New("Data mismatch: Uploader doesn't match provided Data.")
	}

	return upload, nil
}

/**
 * Reconstruct an upload from the tx metadata, ie /tx/<id>.
 *
 * @param api
 * @param id
 * @param data
 */
func (tt *TransactionUploader) FromTransactionId(id string) (*SerializedUploader, error) {
	return nil, nil
}

func (tt *TransactionUploader) FormatSerializedUploader() *SerializedUploader {
	tx := tt.Transaction
	return &SerializedUploader{
		chunkIndex:         tt.ChunkIndex,
		txPosted:           tt.TxPosted,
		transaction:        tx,
		lastRequestTimeEnd: tt.LastRequestTimeEnd,
		lastResponseStatus: tt.LastResponseStatus,
		lastResponseError:  tt.LastResponseError,
	}
}

// POST to /tx
func (tt *TransactionUploader) postTransaction() error {
	var uploadInBody = tt.TotalChunks() <= MAX_CHUNKS_IN_BODY
	return tt.uploadTx(uploadInBody)
}

func (tt *TransactionUploader) uploadTx(withBody bool) error {
	// if withBody {
	// 	// Post the Transaction with Data.
	// 	tt.Transaction.Data = utils.Base64Encode(tt.Data)
	// }
	body, statusCode, err := tt.Client.SubmitTransaction(tt.Transaction)
	if err != nil || statusCode >= 400 {
		tt.LastResponseError = fmt.Sprintf("%v,%s", err, body)
		tt.LastResponseStatus = statusCode
		return errors.New(fmt.Sprintf("Unable to upload Transaction: %d, %v, %s", statusCode, err, body))
	}

	tt.LastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = statusCode

	// if withBody {
	// 	tt.Transaction.Data = ""
	// }

	// tx already processed
	if statusCode >= 200 && statusCode < 300 {
		tt.TxPosted = true
		// if withBody {
		// 	// We are complete.
		// 	tt.ChunkIndex = tx.MAX_CHUNKS_IN_BODY
		// }
		return nil
	}

	// if withBody {
	// 	tt.LastResponseError = ""
	// }
	return nil
}
