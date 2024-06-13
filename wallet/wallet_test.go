package wallet

import (
	"testing"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/transaction"
	"github.com/stretchr/testify/assert"
)

func mint(t *testing.T, c *client.Client, address string) {
	_, err := c.Client.Get(c.Gateway + "/mint/" + address + "/10000")
	assert.NoError(t, err)
	mine(t, c)
}

func mine(t *testing.T, c *client.Client) {
	_, err := c.Client.Get(c.Gateway + "/mine")
	assert.NoError(t, err)
}

func createTransaction(t *testing.T, w *Wallet) *transaction.Transaction {
	data := []byte{1, 2, 3}
	tx := transaction.New(data, nil, "", "0", "0")

	tx.Owner = w.Signer.Owner()

	anchor, err := w.Client.GetTransactionAnchor()
	assert.NoError(t, err)
	tx.LastTx = anchor

	reward, err := w.Client.GetTransactionPrice(len(data), "")
	assert.NoError(t, err)
	tx.Reward = reward

	mint(t, w.Client, w.Signer.Address)

	err = tx.Sign(w.Signer)
	assert.NoError(t, err)
	_, err = w.Client.SubmitTransaction(tx)
	assert.NoError(t, err)
	mine(t, w.Client)

	return tx
}

func TestSignTransaction(t *testing.T) {
	w, err := FromPath("../test/signer.json", "http://localhost:1984")
	assert.NoError(t, err)

	data := []byte{1, 2, 3}

	t.Run("Sign", func(t *testing.T) {

		tx := transaction.New(data, nil, "", "0", "0")
		tx, err = w.SignTransaction(tx)
		assert.NoError(t, err)
		assert.NotEmpty(t, tx.ID)
		assert.NotEmpty(t, tx.Signature)
	})
}

func TestSendTransaction(t *testing.T) {
	w, err := FromPath("../test/signer.json", "http://localhost:1984")
	assert.NoError(t, err)

	tx := createTransaction(t, w)

	t.Run("Sent", func(t *testing.T) {
		err = w.SendTransaction(tx)
		assert.NoError(t, err)
	})

	t.Run("ID or Signature not found", func(t *testing.T) {
		tx := createTransaction(t, w)
		tx.ID = ""
		err = w.SendTransaction(tx)
		assert.Error(t, err)

		tx = createTransaction(t, w)
		tx.Signature = ""
		err = w.SendTransaction(tx)
		assert.Error(t, err)
	})
}
