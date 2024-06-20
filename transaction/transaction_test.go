package transaction

import (
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	data := []byte("test")

	s, err := signer.FromPath("../test/signer.json")
	assert.NoError(t, err)

	t.Run("Sign", func(t *testing.T) {
		tx := New(data, "", "0", nil)
		assert.NoError(t, err)
		tx.Owner = s.Owner()
		tx.LastTx = "lqsw6xgaaunfs8h3d6n54ci1lgm2tmtqvz3wke9v9ygq64q8s68yz2jfq5xy4nec"
		tx.Reward = "1000"

		err = tx.Sign(s)
		assert.NoError(t, err)

		err = tx.Verify()
		assert.NoError(t, err)
	})

	t.Run("Sign with Tags", func(t *testing.T) {
		tags := &[]tag.Tag{{Name: "test", Value: "test"}, {Name: "test", Value: "1"}, {Name: "test", Value: "test"}}
		tx := New(data, "", "0", tags)
		assert.NoError(t, err)
		tx.Owner = s.Owner()
		tx.LastTx = "lqsw6xgaaunfs8h3d6n54ci1lgm2tmtqvz3wke9v9ygq64q8s68yz2jfq5xy4nec"
		tx.Reward = "1000"

		err = tx.Sign(s)
		assert.NoError(t, err)

		err = tx.Verify()
		assert.NoError(t, err)
	})

}
