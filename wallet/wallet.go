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
	tx.Owner = w.Signer.Owner()

	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return err
	}
	tx.LastTx = anchor

	reward, err := w.Client.GetTransactionPrice(len(tx.Data), "")
	if err != nil {
		return err
	}
	tx.Reward = reward
	
	if err = tx.Sign(w.Signer); err != nil {
		return err
	}
	return nil
}

func (w *Wallet) SendTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	if tx.ID == "" || tx.Signature == "" {
		return nil, errors.New("transaction not signed")
	}
	tu, err := uploader.New(w.Client, tx)
	if err != nil {
		return nil, err
	}
	if err = tu.PostTransaction(); err != nil {
		return nil, err
	}
	return tx, nil
}
