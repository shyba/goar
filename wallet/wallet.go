package wallet

import (
	"errors"
	"os"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/transaction"
	"github.com/liteseed/goar/uploader"
)

type Wallet struct {
	Client *client.Client
	Signer *signer.Signer
}

func New(gateway string) (w *Wallet, err error) {
	jwk, err := signer.New()
	if err != nil {
		return nil, err
	}
	return FromJWK(jwk, gateway)
}

func FromPath(path string, gateway string) (*Wallet, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromJWK(b, gateway)
}

func FromJWK(jwk []byte, gateway string) (*Wallet, error) {
	signer, err := signer.FromJWK(jwk)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Client: client.New(gateway),
		Signer: signer,
	}, nil
}

func (w *Wallet) SignTransaction(tx *transaction.Transaction) error {
	anchor, err := w.Client.GetLastTransactionID(w.Signer.Address)
	if err != nil {
		return err
	}
	tx.LastTx = anchor
	tx.Owner = w.Signer.Owner()
	if err = w.Signer.SignTransaction(tx); err != nil {
		return err
	}
	return nil
}

func (w *Wallet) SendTransaction(t *transaction.Transaction) (*transaction.Transaction, error) {
	if t.ID == "" || t.Signature == "" {
		return nil, errors.New("transaction not signed")
	}
	tu, err := uploader.New(w.Client, t)
	if err != nil {
		return nil, err
	}
	if err = tu.PostTransaction(); err != nil {
		return nil, err
	}
	return t, nil
}
