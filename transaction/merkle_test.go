package transaction

import (
	"os"
	"testing"

	"github.com/liteseed/goar/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	root     = "t-GCOnjPWxdox950JsrFMu3nzOE4RktXpMcIlkqSUTw"
	path     = "7EAC9FsACQRwe4oIzu7Mza9KjgWKT4toYxDYGjWrCdp0QgsrYS6AueMJ_rM6ZEGslGqjUekzD3WSe7B5_fwipgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAnH6dASdQCigcL43lp0QclqBaSncF4TspuvxoFbn2L18EXpQrP1wkbwdIjSSWQQRt_F31yNvxtc09KkPFtzMKAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAIHiHU9QwOImFzjqSlfxkJJCtSbAox6TbbFhQvlEapSgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAA"
	offset   = 262143
	dataSize = 836907
)

func TestMerkle(t *testing.T) {
	data, err := os.ReadFile("../test/stubs/rebar3")
	assert.NoError(t, err)

	t.Run("should build a tree with a valid root", func(t *testing.T) {
		rootNode, err := generateTree(data)
		assert.NoError(t, err)
		assert.Equal(t, root, crypto.Base64Encode(rootNode.ID))
	})

	t.Run("should build valid proofs from tree", func(t *testing.T) {
		rootNode, err := generateTree(data)
		assert.NoError(t, err)
		proofs := generateProofs(rootNode, nil, 0)

		assert.Equal(t, path, crypto.Base64Encode(proofs[0].Proof))
	})

	t.Run("should flatten", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5}, flatten[int]([]any{1, []any{2, 3, []any{4, 5}}}))
		assert.Equal(t, []int{1, 2, 3}, flatten[int]([]any{1, []any{2, 3}}))
		assert.Equal(t, []int{1}, flatten[int]([]any{1}))
		assert.Equal(t, []int{1}, flatten[int]([]any{[]any{[]any{1}}}))
	})
}
