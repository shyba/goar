package tx

import (
	"errors"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/types"
)

func GetTransactionDeepHash(t *types.Transaction) ([]byte, error) {
	if t.Format != 2 {
		return nil, errors.New("only type 2 transaction supported")
	}
	rawOwner, err := crypto.Base64Decode(t.Owner)
	if err != nil {
		return nil, err
	}
	rawTarget, err := crypto.Base64Decode(t.Target)
	if err != nil {
		return nil, err
	}

	rawTags, err := DecodeTags(t.Tags)
	if err != nil {
		return nil, err
	}
	rawData, err := crypto.Base64Decode(t.Data)
	if err != nil {
		return nil, err
	}
	rawLastTx, err := crypto.Base64Decode(t.LastTx)
	if err != nil {
		return nil, err
	}
	err = t.PrepareChunks(rawData)
	if err != nil {
		return nil, err
	}

	rawDataRoot, err := crypto.Base64Decode(t.DataRoot)
	if err != nil {
		return nil, err
	}

	chunks := []any{
		"2",
		rawOwner,
		rawTarget,
		rawData,
		[]byte(t.Quantity),
		[]byte(t.Reward),
		rawLastTx,
		rawTags,
		[]byte(t.DataSize),
		rawDataRoot,
	}
	signatureData := crypto.DeepHash(chunks)
	deepHash := signatureData[:]
	return deepHash, nil

}
