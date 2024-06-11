package signer

import (
	"testing"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/transaction"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	c := client.New("http://localhost:1984")
	data := []byte("test")
	// jwk, err := New()
	// assert.NoError(t, err)

	s, err := FromPath("../test/signer.json")
	assert.NoError(t, err)

	tx := transaction.New(data, nil, "", "0", "0")
	assert.NoError(t, err)

	t.Run("Sign", func(t *testing.T) {
		tx.Owner = s.Owner()

		anchor, err := c.GetLastTransactionID(s.Address)
		assert.NoError(t, err)
		tx.LastTx = anchor

		reward, err := c.GetTransactionPrice(len(data), "")
		assert.NoError(t, err)
		tx.Reward = reward

		err = s.SignTransaction(tx)
		assert.NoError(t, err)

		err = transaction.Verify(tx)
		assert.NoError(t, err)
	})
}
