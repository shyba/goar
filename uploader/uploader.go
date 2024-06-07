package uploader

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/types"
)

const (
	DEFAULT_CHUNK_CONCURRENT_NUM = 50
	MAX_CHUNKS_IN_BODY           = 1
	DELAY                        = 30000
)

var FATAL_CHUNK_UPLOAD_ERRORS = []string{
	"invalid_json",
	"chunk_too_big",
	"data_path_too_big",
	"offset_too_big",
	"data_size_too_big",
	"chunk_proof_ratio_not_attractive",
	"invalid_proof",
}

type TransactionUploader struct {
	client             *client.Client
	transaction        *types.Transaction
	ChunkIndex         int
	TxPosted           bool
	Transaction        *types.Transaction
	Data               []byte
	LastRequestTimeEnd int64
	TotalErrors        int // Not serialized.
	LastResponseStatus int
	LastResponseError  string
	TotalChunks        int
}

func New(c *client.Client, t *types.Transaction) (*TransactionUploader, error) {
	return &TransactionUploader{
		client:             c,
		transaction:        t,
		ChunkIndex:         0,
		TxPosted:           false,
		Transaction:        t,
		Data:               nil,
		LastRequestTimeEnd: 0,
		TotalErrors:        0,
		LastResponseStatus: 0,
		LastResponseError:  "",
		TotalChunks:        0,
	}, nil
}

func (tu *TransactionUploader) PostTransaction() error {
	uploadInBody := tu.TotalChunks <= MAX_CHUNKS_IN_BODY
	if uploadInBody {

		_, code, err := tu.client.SubmitTransaction(tu.transaction)
		if err != nil {
			return err
		}
		tu.LastRequestTimeEnd = time.Now().UTC().UnixMilli()
		tu.LastResponseStatus = code
		if code >= 200 && code < 300 {
			tu.TxPosted = true
			tu.ChunkIndex = MAX_CHUNKS_IN_BODY
		}
		return nil
	} else {
		// Post transaction with no data
		t := tu.transaction
		t.Data = ""
		_, code, err := tu.client.SubmitTransaction(t)
		if err != nil {
			return err
		}
		tu.LastRequestTimeEnd = time.Now().UTC().UnixMilli()
		tu.LastResponseStatus = code
		if code >= 200 && code < 300 {
			tu.TxPosted = true
			return nil
		}
		return nil
	}
}

func (tu *TransactionUploader) UploadChunk(chunkIndex int) error {
	if tu.TxPosted && tu.ChunkIndex == len(tu.transaction.ChunkData.Chunks) {
		return errors.New("upload is already complete")
	}

	if tu.LastResponseError != "" {
		tu.TotalErrors++
	} else {
		tu.TotalErrors = 0
	}

	if tu.TotalErrors == 100 {
		return fmt.Errorf("fatal: unable to complete upload: %d: %s", tu.LastResponseStatus, tu.LastResponseError)
	}

	var delay = 0.0
	if tu.LastResponseError != "" {
		delay = DELAY + math.Max(0, float64(tu.LastRequestTimeEnd)-float64(time.Now().UTC().UnixMilli()))
	}

	if delay > 0 {
		delay = delay - delay*0.3*rand.Float64()
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	if !tu.TxPosted {
		return tu.PostTransaction()
	}

	chunk, err := tu.transaction.GetChunk(chunkIndex, tu.Data)
	if err != nil {
		return err
	}

	body, statusCode, err := tu.client.UploadChunk(chunk)
	tu.LastRequestTimeEnd = time.Hour.Milliseconds()
	tu.LastResponseStatus = statusCode

	if tu.LastResponseStatus == 200 {
		tu.ChunkIndex++
	} else {
		if body != nil {
			tu.LastResponseError = string(body)
		} else {
			tu.LastResponseError = err.Error()
		}
		if slices.Contains(FATAL_CHUNK_UPLOAD_ERRORS, tu.LastResponseError) {
			return fmt.Errorf("fatal: unable to complete upload: %d: %s", tu.LastResponseStatus, tu.LastResponseError)
		}
	}
	return nil
}
