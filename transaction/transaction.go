package transaction

import (
	"errors"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tag"
)

func New(data []byte, tags []tag.Tag, target string, quantity string, reward string) *Transaction {
	if tags == nil {
		tags = []tag.Tag{}
	}
	if quantity == "" {
		quantity = "0"
	}
	if reward == "" {
		reward = "0"
	}
	return &Transaction{
		Format:   2,
		Data:     data,
		Target:   target,
		Quantity: quantity,
		Reward:   reward,
		Tags:     tags,
	}
}

func GetTransactionDeepHash(tx *Transaction) ([]byte, error) {
	if tx.Format != 2 {
		return nil, errors.New("only type 2 transaction supported")
	}
	rawOwner, err := crypto.Base64Decode(tx.Owner)
	if err != nil {
		return nil, err
	}
	rawTarget, err := crypto.Base64Decode(tx.Target)
	if err != nil {
		return nil, err
	}

	rawTags, err := tag.Decode(tx.Tags)
	if err != nil {
		return nil, err
	}
	rawLastTx, err := crypto.Base64Decode(tx.LastTx)
	if err != nil {
		return nil, err
	}

	err = tx.PrepareChunks(tx.Data)
	if err != nil {
		return nil, err
	}

	rawDataRoot, err := crypto.Base64Decode(tx.DataRoot)
	if err != nil {
		return nil, err
	}

	chunks := []any{
		[]byte("2"),
		rawOwner,
		rawTarget,
		[]byte(tx.Quantity),
		[]byte(tx.Reward),
		rawLastTx,
		rawTags,
		[]byte(tx.DataSize),
		rawDataRoot,
	}

	deepHash := crypto.DeepHash(chunks)
	signatureData := deepHash[:]
	return signatureData, nil
}

func Verify(tx *Transaction) error {
	signatureData, err := GetTransactionDeepHash(tx)
	if err != nil {
		return err
	}
	rawSignature, err := crypto.Base64Decode(tx.Signature)
	if err != nil {
		return err
	}
	publicKey, err := crypto.GetPublicKeyFromOwner(tx.Owner)
	if err != nil {
		return err
	}
	return crypto.Verify(signatureData, rawSignature, publicKey)
}
