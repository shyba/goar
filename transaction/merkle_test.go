// Package transaction tests - verifies Merkle tree functionality for data chunking and verification
package transaction

import (
	"os"
	"strconv"
	"testing"

	"github.com/liteseed/goar/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test constants for known valid Merkle tree data
const (
	rootBase64URL = "t-GCOnjPWxdox950JsrFMu3nzOE4RktXpMcIlkqSUTw"                                                                                                                                                                                                                                                                                                            // Expected root hash for rebar3 test file
	pathBase64URL = "7EAC9FsACQRwe4oIzu7Mza9KjgWKT4toYxDYGjWrCdp0QgsrYS6AueMJ_rM6ZEGslGqjUekzD3WSe7B5_fwipgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAnH6dASdQCigcL43lp0QclqBaSncF4TspuvxoFbn2L18EXpQrP1wkbwdIjSSWQQRt_F31yNvxtc09KkPFtzMKAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAIHiHU9QwOImFzjqSlfxkJJCtSbAox6TbbFhQvlEapSgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAA" // Expected proof path for first chunk
	offset        = 262143                                                                                                                                                                                                                                                                                                                                                   // Expected offset for test data
	dataSize      = 836907                                                                                                                                                                                                                                                                                                                                                   // Expected data size for test data
)

// TestMerkle verifies comprehensive Merkle tree functionality
func TestMerkle(t *testing.T) {
	t.Run("should validate all paths in 1MB.bin test file", func(t *testing.T) {
		// Load test data file
		data, err := os.ReadFile("../test/1MB.bin")
		require.NoError(t, err)

		// Create transaction and prepare chunks
		tx := New(data, "", "", nil)
		tx.LastTx = "foo"
		tx.Reward = "1"

		err = tx.PrepareChunks(data)
		require.NoError(t, err)
		require.NotNil(t, tx.ChunkData)

		// Decode the data root for validation
		txDataRoot, err := crypto.Base64URLDecode(tx.DataRoot)
		require.NoError(t, err)

		// Validate each chunk's Merkle proof
		for i := range tx.ChunkData.Chunks {
			chunk, err := tx.GetChunk(i, data)
			require.NoError(t, err)

			offset, err := strconv.Atoi(chunk.Offset)
			require.NoError(t, err)

			dataSize, err := strconv.Atoi(chunk.DataSize)
			require.NoError(t, err)

			dataPath, err := crypto.Base64URLDecode(chunk.DataPath)
			require.NoError(t, err)

			// Validate that the chunk belongs to the tree
			result, err := validatePath(txDataRoot, offset, 0, dataSize, dataPath)
			assert.NotNil(t, result)
			assert.NoError(t, err)
		}
	})

	t.Run("should validate all paths in lotsofdata.bin test file", func(t *testing.T) {
		// Load larger test data file
		data, err := os.ReadFile("../test/lotsofdata.bin")
		require.NoError(t, err)

		// Create transaction and prepare chunks
		tx := New(data, "", "", nil)
		tx.LastTx = "foo"
		tx.Reward = "1"

		err = tx.PrepareChunks(data)
		require.NoError(t, err)
		require.NotNil(t, tx.ChunkData)

		// Decode the data root for validation
		txDataRoot, err := crypto.Base64URLDecode(tx.DataRoot)
		require.NoError(t, err)

		// Validate each chunk's Merkle proof
		for i := range tx.ChunkData.Chunks {
			chunk, err := tx.GetChunk(i, data)
			require.NoError(t, err)

			offset, err := strconv.Atoi(chunk.Offset)
			require.NoError(t, err)

			dataSize, err := strconv.Atoi(chunk.DataSize)
			require.NoError(t, err)

			dataPath, err := crypto.Base64URLDecode(chunk.DataPath)
			require.NoError(t, err)

			// Validate that the chunk belongs to the tree
			result, err := validatePath(txDataRoot, offset, 0, dataSize, dataPath)
			assert.NotNil(t, result)
			assert.NoError(t, err)
		}
	})

	t.Run("should build a tree with a valid root", func(t *testing.T) {
		// Load test file with known expected root
		data, err := os.ReadFile("../test/rebar3")
		require.NoError(t, err)

		// Generate Merkle tree
		rootNode, err := generateTree(data)
		require.NoError(t, err)
		require.NotNil(t, rootNode)

		// Verify the root hash matches expected value
		assert.Equal(t, rootBase64URL, crypto.Base64URLEncode(rootNode.ID))
	})

	t.Run("should build valid proofs from tree", func(t *testing.T) {
		// Load test file
		data, err := os.ReadFile("../test/rebar3")
		require.NoError(t, err)

		// Generate Merkle tree
		rootNode, err := generateTree(data)
		require.NoError(t, err)

		// Generate proofs for all chunks
		proofs := generateProofs(rootNode, nil, 0)
		require.NotEmpty(t, proofs)

		// Verify the first proof matches expected value
		assert.Equal(t, pathBase64URL, crypto.Base64URLEncode(proofs[0].Proof))
	})

	t.Run("should flatten nested slices correctly", func(t *testing.T) {
		// Test the flatten utility function with various nested structures
		assert.Equal(t, []int{1, 2, 3, 4, 5}, flatten[int]([]any{1, []any{2, 3, []any{4, 5}}}))
		assert.Equal(t, []int{1, 2, 3}, flatten[int]([]any{1, []any{2, 3}}))
		assert.Equal(t, []int{1}, flatten[int]([]any{1}))
		assert.Equal(t, []int{1}, flatten[int]([]any{[]any{[]any{1}}}))
	})

	t.Run("should reject invalid Merkle path", func(t *testing.T) {
		// Create an invalid proof path
		invalidPath, err := crypto.Base64URLDecode(
			"VUSdubFW2cTvvr5s6VGSU2oxftxma77bRvils5fqikdj4qnP8xEG2HQQKyZeZGW5b9WNFlmDRBTyTJ8NnHQD3tLHc2VwctfdrXbkUODANATrOP6p8RNlSNT50jMKdSKymG0M8yv9g3LCoPB4QXawcRP6q9X5u1nnI7GFMlyuxoC4p21zWi7v68f1r73wXHWdH76VgCNbt0lEUDg1pW8sYvi6pdwAdTNdQIcAhqkO2JBJ2Kwtlxemj4E6NMKg9wi2pQHt6CKlX3T5rQdVd0Tt8czxrkOUBAW9J8XGK9iSLoj4LWZl3z4cKIFyZH7iUgIzCu9Id8jIoO93lVdgaUa4RW",
		)
		require.NoError(t, err)

		// Get the known good root hash
		root, err := crypto.Base64URLDecode(rootBase64URL)
		require.NoError(t, err)

		// Attempt to validate the invalid path - should fail
		result, err := validatePath(root, offset, 0, dataSize, invalidPath)
		assert.Nil(t, result)
		assert.Error(t, err)
	})
}
