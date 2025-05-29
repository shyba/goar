// Package transaction tests - verifies transaction creation, signing, and verification
package transaction

import (
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSign verifies transaction signing and verification functionality
func TestSign(t *testing.T) {
	data := []byte("test")

	s, err := signer.FromPath("../test/signer.json")
	require.NoError(t, err)

	t.Run("Sign basic transaction", func(t *testing.T) {
		tx := New(data, "", "0", nil)
		require.NotNil(t, tx)

		// Set required fields for signing
		tx.Owner = s.Owner()
		tx.LastTx = "lqsw6xgaaunfs8h3d6n54ci1lgm2tmtqvz3wke9v9ygq64q8s68yz2jfq5xy4nec"
		tx.Reward = "1000"

		// Sign the transaction
		err = tx.Sign(s)
		require.NoError(t, err)

		// Verify signature was created
		assert.NotEmpty(t, tx.ID)
		assert.NotEmpty(t, tx.Signature)

		// Verify the signature is valid
		err = tx.Verify()
		assert.NoError(t, err)
	})

	t.Run("Sign transaction with tags", func(t *testing.T) {
		tags := &[]tag.Tag{
			{Name: "test", Value: "test"},
			{Name: "test", Value: "1"},
			{Name: "test", Value: "test"},
		}
		tx := New(data, "", "0", tags)
		require.NotNil(t, tx)

		// Set required fields for signing
		tx.Owner = s.Owner()
		tx.LastTx = "lqsw6xgaaunfs8h3d6n54ci1lgm2tmtqvz3wke9v9ygq64q8s68yz2jfq5xy4nec"
		tx.Reward = "1000"

		// Sign the transaction
		err = tx.Sign(s)
		require.NoError(t, err)

		// Verify signature was created
		assert.NotEmpty(t, tx.ID)
		assert.NotEmpty(t, tx.Signature)
		assert.NotNil(t, tx.Tags)
		assert.Len(t, *tx.Tags, 3) // Should have 3 tags

		// Verify the signature is valid
		err = tx.Verify()
		assert.NoError(t, err)
	})
}

// TestNew verifies transaction creation with various parameters
func TestNew(t *testing.T) {
	t.Run("Create transaction with data", func(t *testing.T) {
		data := []byte("hello world")
		tx := New(data, "", "0", nil)

		assert.Equal(t, 2, tx.Format)
		assert.NotEmpty(t, tx.Data)
		assert.Equal(t, "", tx.Target)
		assert.Equal(t, "0", tx.Quantity)
		assert.NotNil(t, tx.Tags)
		assert.Equal(t, "0", tx.DataSize)
	})

	t.Run("Create AR transfer transaction", func(t *testing.T) {
		target := "test_address"
		quantity := "1000000000000" // 1 AR in Winston
		tx := New(nil, target, quantity, nil)

		assert.Equal(t, 2, tx.Format)
		assert.Equal(t, target, tx.Target)
		assert.Equal(t, quantity, tx.Quantity)
		assert.NotNil(t, tx.Tags)
	})

	t.Run("Create transaction with tags", func(t *testing.T) {
		tags := &[]tag.Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "App-Name", Value: "Test-App"},
		}
		tx := New([]byte("test"), "", "0", tags)

		assert.NotNil(t, tx.Tags)
		assert.Len(t, *tx.Tags, 2) // Should have 2 tags
		// Note: New() converts tags to base64url format, so we can't directly compare
	})
}
