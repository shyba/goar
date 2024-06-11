package client

import (
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/transaction"
	"github.com/stretchr/testify/assert"
)

// func TestGetTransactionByID(t *testing.T) {
// 	c := New("http://localhost:1984")

// 	jwk, err := signer.New()
// 	assert.NoError(t, err)

// 	s, err := signer.FromJWK(jwk)
// 	assert.NoError(t, err)

// 	transaction := &types.Transaction{Format: 2, Data: crypto.Base64Encode([]byte{1, 2, 3})}
// 	err = s.SignTransaction(transaction)

// 	assert.NoError(t, err)

// 	t.Run("found", func(t *testing.T) {
// 		_, _, _ = c.SubmitTransaction(transaction)
// 		f, err := c.GetTransactionByID(transaction.ID)
// 		assert.NoError(t, err)
// 		assert.Equal(t, transaction, f)
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		f, err := c.GetTransactionByID(transaction.ID)
// 		assert.Nil(t, f)
// 		assert.Error(t, errors.New("Not Found"), err)
// 	})
// }

func TestSubmitTransaction(t *testing.T) {
	c := New("http://localhost:1984")
	data := []byte("test")
	jwk, err := signer.New()
	assert.NoError(t, err)

	s, err := signer.FromJWK(jwk)
	assert.NoError(t, err)

	tx := transaction.New(data, nil, "", "0", "0")
	assert.NoError(t, err)

	tx.Owner = s.Owner()

	anchor, err := c.GetLastTransactionID(s.Address)
	assert.NoError(t, err)
	tx.LastTx = anchor

	reward, err := c.GetTransactionPrice(len(data), "")
	assert.NoError(t, err)
	tx.Reward = reward

	err = s.SignTransaction(tx)
	assert.NoError(t, err)

	t.Run("Post", func(t *testing.T) {
		res, code, err := c.SubmitTransaction(tx)
		assert.Equal(t, "ok", string(res))
		assert.Equal(t, 201, code)
		assert.NoError(t, err)
	})
}
