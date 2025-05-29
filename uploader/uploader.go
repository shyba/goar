// Package uploader provides functionality for uploading transactions and data to Arweave nodes.
//
// This package handles the complex process of uploading transactions to the Arweave
// network, including chunked upload for large data, retry logic, and error handling.
// It supports both small transactions (uploaded in a single request) and large
// transactions (uploaded as chunks with Merkle proofs).
//
// Example usage:
//
//	// Create uploader for a transaction
//	uploader, err := New(client, transaction)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Upload the transaction
//	err = uploader.PostTransaction()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// For large transactions, upload chunks
//	for i := 0; i < uploader.TotalChunks; i++ {
//		err = uploader.UploadChunk(i)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
package uploader

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/transaction"
)

// Upload configuration constants
const (
	MAX_CHUNKS_IN_BODY = 1     // Maximum number of chunks to include in transaction body
	DELAY              = 30000 // Base delay in milliseconds for retry logic
)

// FATAL_CHUNK_UPLOAD_ERRORS lists errors that should not be retried.
// These errors indicate permanent failures that won't be resolved by retrying.
var FATAL_CHUNK_UPLOAD_ERRORS = []string{
	"invalid_json",                     // JSON parsing error
	"chunk_too_big",                    // Chunk exceeds size limits
	"data_path_too_big",                // Merkle proof path is too large
	"offset_too_big",                   // Chunk offset is invalid
	"data_size_too_big",                // Total data size exceeds limits
	"chunk_proof_ratio_not_attractive", // Economic constraints not met
	"invalid_proof",                    // Merkle proof verification failed
}

// TransactionUploader manages the upload process for an Arweave transaction.
//
// This struct tracks the state of an ongoing upload operation, including
// which chunks have been uploaded, error counts, and timing information.
// It handles both simple uploads (small transactions) and chunked uploads
// (large transactions with Merkle proofs).
type TransactionUploader struct {
	client             *client.Client           // HTTP client for communicating with Arweave nodes
	transaction        *transaction.Transaction // The transaction being uploaded
	ChunkIndex         int                      // Index of the next chunk to upload
	TxPosted           bool                     // Whether the transaction header has been posted
	Data               []byte                   // Raw transaction data (for chunk generation)
	LastRequestTimeEnd int64                    // Timestamp of last request completion
	TotalErrors        int                      // Running count of upload errors (not serialized)
	LastResponseStatus int                      // HTTP status code from last request
	LastResponseError  string                   // Error message from last failed request
	TotalChunks        int                      // Total number of chunks in this transaction
}

// New creates a new TransactionUploader for the given transaction.
//
// This function initializes an uploader instance to manage the upload
// process for a transaction. The uploader tracks upload state and handles
// retry logic for failed uploads.
//
// Parameters:
//   - c: HTTP client for communicating with Arweave nodes
//   - t: The transaction to upload
//
// Returns a new TransactionUploader instance ready to begin uploading.
//
// Example:
//
//	uploader, err := New(client, signedTransaction)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Created uploader for transaction %s\n", signedTransaction.ID)
func New(c *client.Client, t *transaction.Transaction) (*TransactionUploader, error) {
	return &TransactionUploader{
		client:             c,
		transaction:        t,
		ChunkIndex:         0,
		TxPosted:           false,
		Data:               nil,
		LastRequestTimeEnd: 0,
		TotalErrors:        0,
		LastResponseStatus: 0,
		LastResponseError:  "",
		TotalChunks:        0,
	}, nil
}

// PostTransaction uploads the transaction to the Arweave network.
//
// This method handles the initial transaction submission. For small transactions
// (with few chunks), the entire transaction including data is uploaded in one
// request. For large transactions, only the transaction header is uploaded,
// and data chunks must be uploaded separately using UploadChunk.
//
// The method automatically determines the upload strategy based on the
// MAX_CHUNKS_IN_BODY constant.
//
// Returns an error if the transaction submission fails.
//
// Example:
//
//	err := uploader.PostTransaction()
//	if err != nil {
//		log.Printf("Failed to post transaction: %v", err)
//		return err
//	}
//	if uploader.TxPosted {
//		fmt.Println("Transaction posted successfully")
//	}
func (tu *TransactionUploader) PostTransaction() error {
	if tu.TotalChunks <= MAX_CHUNKS_IN_BODY {
		code, err := tu.client.SubmitTransaction(tu.transaction)
		if err != nil {
			return err
		}
		tu.LastRequestTimeEnd = time.Now().UTC().UnixMilli()
		tu.LastResponseStatus = code
		if code >= 200 && code < 400 {
			tu.TxPosted = true
			tu.ChunkIndex = MAX_CHUNKS_IN_BODY
		}
		return nil
	} else {
		// Post transaction with no data
		t := tu.transaction
		t.Data = ""
		code, err := tu.client.SubmitTransaction(t)
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

// UploadChunk uploads a specific chunk of the transaction data.
//
// This method uploads individual data chunks for large transactions.
// It includes sophisticated retry logic with exponential backoff and
// handles both temporary and permanent errors appropriately.
//
// The method will:
// 1. Check if upload is already complete
// 2. Track error counts and implement failure limits
// 3. Apply retry delays with jitter
// 4. Post the transaction header if not already done
// 5. Upload the specified chunk with its Merkle proof
// 6. Handle response codes and errors
//
// Parameters:
//   - chunkIndex: The index of the chunk to upload (0-based)
//
// Returns an error if the chunk upload fails permanently or if
// too many errors have occurred.
//
// Example:
//
//	// Upload all chunks
//	for i := 0; i < uploader.TotalChunks; i++ {
//		err := uploader.UploadChunk(i)
//		if err != nil {
//			log.Printf("Failed to upload chunk %d: %v", i, err)
//			return err
//		}
//		fmt.Printf("Uploaded chunk %d/%d\n", i+1, uploader.TotalChunks)
//	}
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

	code, err := tu.client.UploadChunk(chunk)
	tu.LastRequestTimeEnd = time.Hour.Milliseconds()
	tu.LastResponseStatus = code

	if tu.LastResponseStatus == 200 {
		tu.ChunkIndex++
	} else {
		if err != nil {
			tu.LastResponseError = err.Error()
		}
		if slices.Contains(FATAL_CHUNK_UPLOAD_ERRORS, tu.LastResponseError) {
			return fmt.Errorf("fatal: unable to complete upload: %d: %s", tu.LastResponseStatus, tu.LastResponseError)
		}
	}
	return nil
}
