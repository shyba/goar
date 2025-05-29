// Package uploader tests - verifies transaction upload functionality
package uploader

import (
	"testing"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew verifies that uploaders can be created correctly
func TestNew(t *testing.T) {
	// Create a mock client and transaction
	client := client.New("http://localhost:1984")
	data := []byte("test data")
	tx := transaction.New(data, "", "0", nil)

	uploader, err := New(client, tx)
	require.NoError(t, err)
	assert.NotNil(t, uploader)
	assert.Equal(t, client, uploader.client)
	assert.Equal(t, tx, uploader.transaction)
	assert.Equal(t, 0, uploader.ChunkIndex)
	assert.False(t, uploader.TxPosted)
	assert.Equal(t, 0, uploader.TotalErrors)
	assert.Equal(t, 0, uploader.LastResponseStatus)
	assert.Empty(t, uploader.LastResponseError)
}

// TestUploaderInitialization verifies uploader is properly initialized
func TestUploaderInitialization(t *testing.T) {
	client := client.New("http://localhost:1984")

	t.Run("Small transaction", func(t *testing.T) {
		data := []byte("small data")
		tx := transaction.New(data, "", "0", nil)

		uploader, err := New(client, tx)
		require.NoError(t, err)
		assert.Equal(t, 0, uploader.TotalChunks)
		assert.Equal(t, 0, uploader.ChunkIndex)
	})

	t.Run("Empty transaction", func(t *testing.T) {
		tx := transaction.New(nil, "target", "1000", nil)

		uploader, err := New(client, tx)
		require.NoError(t, err)
		assert.NotNil(t, uploader)
	})
}

// TestFatalErrors verifies fatal error detection
func TestFatalErrors(t *testing.T) {
	testCases := []struct {
		name    string
		error   string
		isFatal bool
	}{
		{"Invalid JSON", "invalid_json", true},
		{"Chunk too big", "chunk_too_big", true},
		{"Data path too big", "data_path_too_big", true},
		{"Offset too big", "offset_too_big", true},
		{"Data size too big", "data_size_too_big", true},
		{"Proof ratio not attractive", "chunk_proof_ratio_not_attractive", true},
		{"Invalid proof", "invalid_proof", true},
		{"Network error", "network_timeout", false},
		{"Temporary error", "temporary_failure", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isFatal := false
			for _, fatalError := range FATAL_CHUNK_UPLOAD_ERRORS {
				if fatalError == tc.error {
					isFatal = true
					break
				}
			}
			assert.Equal(t, tc.isFatal, isFatal)
		})
	}
}

// TestConstants verifies important constants are set correctly
func TestConstants(t *testing.T) {
	assert.Equal(t, 1, MAX_CHUNKS_IN_BODY)
	assert.Equal(t, 30000, DELAY)
	assert.Len(t, FATAL_CHUNK_UPLOAD_ERRORS, 7)
}

// TestUploaderFields verifies all uploader fields are accessible
func TestUploaderFields(t *testing.T) {
	client := client.New("http://localhost:1984")
	data := []byte("test data for uploader")
	tx := transaction.New(data, "", "0", nil)

	uploader, err := New(client, tx)
	require.NoError(t, err)

	// Test that we can access and modify all fields
	uploader.ChunkIndex = 5
	uploader.TxPosted = true
	uploader.Data = []byte("new data")
	uploader.LastRequestTimeEnd = 123456789
	uploader.TotalErrors = 3
	uploader.LastResponseStatus = 200
	uploader.LastResponseError = "test error"
	uploader.TotalChunks = 10

	assert.Equal(t, 5, uploader.ChunkIndex)
	assert.True(t, uploader.TxPosted)
	assert.Equal(t, []byte("new data"), uploader.Data)
	assert.Equal(t, int64(123456789), uploader.LastRequestTimeEnd)
	assert.Equal(t, 3, uploader.TotalErrors)
	assert.Equal(t, 200, uploader.LastResponseStatus)
	assert.Equal(t, "test error", uploader.LastResponseError)
	assert.Equal(t, 10, uploader.TotalChunks)
}

// MockTransaction creates a properly signed transaction for testing
func createMockSignedTransaction(t *testing.T) *transaction.Transaction {
	s, err := signer.FromPath("../test/signer.json")
	require.NoError(t, err)

	data := []byte("test transaction data")
	tx := transaction.New(data, "", "0", nil)
	tx.Owner = s.Owner()
	tx.LastTx = "test_anchor"
	tx.Reward = "1000"

	err = tx.Sign(s)
	require.NoError(t, err)

	return tx
}

// TestPostTransactionValidation verifies transaction validation before posting
func TestPostTransactionValidation(t *testing.T) {
	client := client.New("http://localhost:1984")
	tx := createMockSignedTransaction(t)

	uploader, err := New(client, tx)
	require.NoError(t, err)

	assert.NotNil(t, uploader.transaction)
	assert.NotEmpty(t, uploader.transaction.ID)
	assert.NotEmpty(t, uploader.transaction.Signature)
}

// Note: Network-dependent tests are commented out as they require a running Arweave node
// These would test the actual upload functionality but need proper test infrastructure

/*
// TestPostTransactionSmall tests posting small transactions
func TestPostTransactionSmall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := client.New("http://localhost:1984")
	tx := createMockSignedTransaction(t)

	uploader, err := New(client, tx)
	require.NoError(t, err)

	// This would require a running Arweave node
	err = uploader.PostTransaction()
	// We can't assert success without a real node, but we can verify the method exists
	assert.NotPanics(t, func() { uploader.PostTransaction() })
}
*/
