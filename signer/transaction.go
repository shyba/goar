package signer

import (
	"crypto/sha256"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tx"
)

func (s *Signer) SignTransaction(t *tx.Transaction) error {
	signatureData, err := tx.GetTransactionDeepHash(t)
	if err != nil {
		return err
	}
	rawSignature, err := crypto.Sign(signatureData, s.PrivateKey)
	if err != nil {
		return err
	}

	txId := sha256.Sum256(rawSignature)
	t.ID = crypto.Base64Encode(txId[:])
	t.Signature = crypto.Base64Encode(rawSignature)
	return nil
}
