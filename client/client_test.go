package client

import (
	"errors"
	"strconv"
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction"
	"github.com/stretchr/testify/assert"
)

func mint(t *testing.T, c *Client, address string) {
	res, err := c.get("mint/" + address + "/1000000000000")
	if err != nil {
		panic(0)
	}
	t.Logf("Balance: %s", string(res))
	mine(c)
}

func mine(c *Client) {
	_, err := c.get("mine")
	if err != nil {
		panic(0)
	}
}

func createTransaction(t *testing.T, c *Client) *transaction.Transaction {
	s, err := signer.FromPath("../test/signer.json")
	assert.NoError(t, err)
	data := []byte{1, 2, 3}

	mint(t, c, s.Address)

	tx := transaction.New(data, "", "0", nil)
	assert.NoError(t, err)

	tx.Owner = s.Owner()

	anchor, err := c.GetTransactionAnchor()
	assert.NoError(t, err)
	tx.LastTx = anchor

	reward, err := c.GetTransactionPrice(len(data), "")
	assert.NoError(t, err)
	tx.Reward = reward

	err = tx.Sign(s)
	assert.NoError(t, err)
	_, err = c.SubmitTransaction(tx)
	assert.NoError(t, err)
	mine(c)

	return tx
}

func TestGetTransactionByID(t *testing.T) {
	c := New("http://localhost:1984")
	tx := createTransaction(t, c)
	t.Run("found", func(t *testing.T) {
		f, err := c.GetTransactionByID(tx.ID)
		assert.NoError(t, err)
		assert.Equal(t, tx.Signature, f.Signature)
	})

	t.Run("not found", func(t *testing.T) {
		f, err := c.GetTransactionByID("QWrt4e6nXe7zNcXJE0IADPZI7f9-O_enUk5g8FE_RpL")
		assert.Nil(t, f)
		assert.Error(t, errors.New("not found"), err)
	})
}

func TestGetTransactionStatus(t *testing.T) {
	c := New("http://localhost:1984")
	tx := createTransaction(t, c)
	_, err := c.GetTransactionStatus(tx.ID)
	assert.NoError(t, err)
}

func TestGetTransactionField(t *testing.T) {
	c := New("http://localhost:1984")
	tx := createTransaction(t, c)
	res, err := c.GetTransactionField(tx.ID, "owner")
	assert.NoError(t, err)
	assert.Equal(t, tx.Owner, res)
}

func TestGetTransactionData(t *testing.T) {
	c := New("http://localhost:1984")
	tx := createTransaction(t, c)
	res, err := c.GetTransactionData(tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, tx.Data, res)
}

func TestGetTransactionPrice(t *testing.T) {
	c := New("http://localhost:1984")
	res, err := c.GetTransactionPrice(0, "")
	assert.NoError(t, err)
	_, err = strconv.Atoi(res)
	assert.NoError(t, err)
}

func TestGetTransactionAnchor(t *testing.T) {
	c := New("http://localhost:1984")
	res, err := c.GetTransactionAnchor()
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestSubmitTransaction(t *testing.T) {
	c := New("http://localhost:1984")
	data := []byte("test")
	tags := &[]tag.Tag{{Name: "test", Value: "test"}, {Name: "test", Value: "test"}, {Name: "test", Value: "test"}}

	s, err := signer.FromPath("../test/signer.json")
	assert.NoError(t, err)

	mine(c)
	mint(t, c, s.Address)

	t.Run("Post with Data", func(t *testing.T) {
		tx := transaction.New(data, "", "0", nil)
		assert.NoError(t, err)

		tx.Owner = s.Owner()

		anchor, err := c.GetTransactionAnchor()
		assert.NoError(t, err)
		tx.LastTx = anchor

		reward, err := c.GetTransactionPrice(len(data), "")
		assert.NoError(t, err)
		tx.Reward = reward

		err = tx.Sign(s)
		assert.NoError(t, err)
		code, err := c.SubmitTransaction(tx)
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	})

	t.Run("Post with Data & Tags", func(t *testing.T) {
		tx := transaction.New(data, "", "0", tags)
		assert.NoError(t, err)

		tx.Owner = s.Owner()

		anchor, err := c.GetTransactionAnchor()
		assert.NoError(t, err)
		tx.LastTx = anchor

		reward, err := c.GetTransactionPrice(len(data), "")
		assert.NoError(t, err)
		tx.Reward = reward

		err = tx.Sign(s)
		assert.NoError(t, err)
		code, err := c.SubmitTransaction(tx)
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	})

	t.Run("Post with Target & Quantity", func(t *testing.T) {
		tx := transaction.New(nil, "Cbj95zDZBBhmyht6iFlEf7xmSCSVZGw436V6HWmm9Ek", "1000", nil)
		assert.NoError(t, err)

		tx.Owner = s.Owner()

		anchor, err := c.GetTransactionAnchor()
		assert.NoError(t, err)
		tx.LastTx = anchor

		reward, err := c.GetTransactionPrice(len(data), "")
		assert.NoError(t, err)
		tx.Reward = reward

		err = tx.Sign(s)
		assert.NoError(t, err)
		code, err := c.SubmitTransaction(tx)
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	})

	t.Run("Post with Target, Quantity, & Tags", func(t *testing.T) {
		tx := transaction.New(nil, "Cbj95zDZBBhmyht6iFlEf7xmSCSVZGw436V6HWmm9Ek", "1000", tags)
		assert.NoError(t, err)

		tx.Owner = s.Owner()

		anchor, err := c.GetTransactionAnchor()
		assert.NoError(t, err)
		tx.LastTx = anchor

		reward, err := c.GetTransactionPrice(len(data), "")
		assert.NoError(t, err)
		tx.Reward = reward

		err = tx.Sign(s)
		assert.NoError(t, err)
		code, err := c.SubmitTransaction(tx)
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	})

	t.Run("Post with Everything", func(t *testing.T) {
		tx := transaction.New(data, "Cbj95zDZBBhmyht6iFlEf7xmSCSVZGw436V6HWmm9Ek", "1000", tags)
		assert.NoError(t, err)

		tx.Owner = s.Owner()

		anchor, err := c.GetTransactionAnchor()
		assert.NoError(t, err)
		tx.LastTx = anchor

		reward, err := c.GetTransactionPrice(len(data), "")
		assert.NoError(t, err)
		tx.Reward = reward

		err = tx.Sign(s)
		assert.NoError(t, err)
		code, err := c.SubmitTransaction(tx)
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	})
}
