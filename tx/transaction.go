package tx

import (
	"fmt"

	"github.com/liteseed/goar/crypto"
)

func GetTransactionChunks(t *Transaction) ([]byte, error) {
	switch t.Format {
	case 2:
		err := t.PrepareChunks(t.Data)
		if err != nil {
			return nil, err
		}
		tags := [][]string{}
		for _, tag := range t.Tags {
			tags = append(tags, []string{
				tag.Name, tag.Value,
			})
		}

		chunks := []any{}

		signatureData := crypto.DeepHash(chunks)
		deepHash := signatureData[:]
		return deepHash, nil

	default:
		return nil, fmt.Errorf("unexpected transaction format: %d", t.Format)
	}
}

// Note: we *do not* use `t.Data`, the caller may be
// operating on a transaction with an zero length data field.
// This function computes the chunks for the data passed in and
// assigns the result to this transaction. It should not read the
// data *from* this transaction.
func (t *Transaction) PrepareChunks(data []byte) error {
	return nil
}
