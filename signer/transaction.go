package signer

import (
	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/transaction"
)

func (s *Signer) SignTransaction(tx *transaction.Transaction) error {
	signatureData, err := transaction.GetTransactionDeepHash(tx)
	if err != nil {
		return err
	}
	rawSignature, err := s.Sign(signatureData)
	if err != nil {
		return err
	}
	txId, err := crypto.SHA256(rawSignature)
	if err != nil {
		return err
	}

	tx.ID = crypto.Base64Encode(txId[:])
	tx.Signature = crypto.Base64Encode(rawSignature)
	return nil
}
