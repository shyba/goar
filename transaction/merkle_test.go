package transaction

import (
	"os"
	"strconv"
	"testing"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/signer"
	"github.com/stretchr/testify/assert"
)

const (
	rootBase64URL = "t-GCOnjPWxdox950JsrFMu3nzOE4RktXpMcIlkqSUTw"
	pathBase64URL = "7EAC9FsACQRwe4oIzu7Mza9KjgWKT4toYxDYGjWrCdp0QgsrYS6AueMJ_rM6ZEGslGqjUekzD3WSe7B5_fwipgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAnH6dASdQCigcL43lp0QclqBaSncF4TspuvxoFbn2L18EXpQrP1wkbwdIjSSWQQRt_F31yNvxtc09KkPFtzMKAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAIHiHU9QwOImFzjqSlfxkJJCtSbAox6TbbFhQvlEapSgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAA"
	offset        = 262143
	dataSize      = 836907
)

func TestMerkle(t *testing.T) {
	t.Run("should validate all paths in 1MB.bin test file", func(t *testing.T) {
		data, err := os.ReadFile("../test/1MB.bin")
		assert.NoError(t, err)

		s, err := signer.FromPath("../test/signer.json")
		assert.NoError(t, err)

		tx := New(data, "", "", nil)
		tx.LastTx = "foo"
		tx.Reward = "1"

		err = tx.Sign(s)
		assert.NoError(t, err)

		err = tx.PrepareChunks(data)
		assert.NoError(t, err)

		txDataRoot, err := crypto.Base64URLDecode(tx.DataRoot)
		assert.NoError(t, err)

		for i := range tx.ChunkData.Chunks {
			chunk, err := tx.GetChunk(i, data)
			assert.NoError(t, err)

			offset, err := strconv.Atoi(chunk.Offset)
			assert.NoError(t, err)

			dataSize, err := strconv.Atoi(chunk.DataSize)
			assert.NoError(t, err)

			dataRoot, err := crypto.Base64URLDecode(chunk.DataRoot)
			assert.NoError(t, err)

			r, err := validatePath(txDataRoot, offset, 0, dataSize, dataRoot)
			assert.NotNil(t, r)
			assert.NoError(t, err)
		}
	})

	t.Run("should build a tree with a valid root", func(t *testing.T) {
		data, err := os.ReadFile("../test/rebar3")
		assert.NoError(t, err)
		rootNode, err := generateTree(data)
		assert.NoError(t, err)
		assert.Equal(t, rootBase64URL, crypto.Base64URLEncode(rootNode.ID))
	})

	t.Run("should build valid proofs from tree", func(t *testing.T) {
		data, err := os.ReadFile("../test/rebar3")
		assert.NoError(t, err)
		rootNode, err := generateTree(data)
		assert.NoError(t, err)
		proofs := generateProofs(rootNode, nil, 0)

		assert.Equal(t, pathBase64URL, crypto.Base64URLEncode(proofs[0].Proof))
	})

	t.Run("should flatten", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5}, flatten[int]([]any{1, []any{2, 3, []any{4, 5}}}))
		assert.Equal(t, []int{1, 2, 3}, flatten[int]([]any{1, []any{2, 3}}))
		assert.Equal(t, []int{1}, flatten[int]([]any{1}))
		assert.Equal(t, []int{1}, flatten[int]([]any{[]any{[]any{1}}}))
	})

	t.Run("should reject invalid path", func(t *testing.T) {
		invalidPath, err := crypto.Base64URLDecode(
			"VUSdubFW2cTvvr5s6VGSU2oxftxma77bRvils5fqikdj4qnP8xEG2HQQKyZeZGW5b9WNFlmDRBTyTJ8NnHQD3tLHc2VwctfdrXbkUODANATrOP6p8RNlSNT50jMKdSKymG0M8yv9g3LCoPB4QXawcRP6q9X5u1nnI7GFMlyuxoC4p21zWi7v68f1r73wXHWdH76VgCNbt0lEUDg1pW8sYvi6pdwAdTNdQIcAhqkO2JBJ2Kwtlxemj4E6NMKg9wi2pQHt6CKlX3T5rQdVd0Tt8czxrkOUBAW9J8XGK9iSLoj4LWZl3z4cKIFyZH7iUgIzCu9Id8jIoO93lVdgaUa4RW",
		)
		assert.NoError(t, err)

		root, err := crypto.Base64URLDecode(rootBase64URL)
		assert.NoError(t, err)

		p, err := validatePath(root, offset, 0, dataSize, invalidPath)
		assert.Nil(t, p)
		assert.Error(t, err)
	})
}
