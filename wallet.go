package goar

import (
	"errors"
	"os"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/types"
	"github.com/liteseed/goar/uploader"
)

type Wallet struct {
	Client *client.Client
	Signer *signer.Signer
}

func New(b []byte, url string) (w *Wallet, err error) {
	signer, err := signer.New(b)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Client: client.New(url),
		Signer: signer,
	}, nil
}

func FromPath(path string, node string) (*Wallet, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return New(b, node)
}

func (w *Wallet) SignTransaction(t *types.Transaction) (*types.Transaction, error) {
	anchor, err := w.Client.GetLastTransactionID(w.Signer.Address)
	if err != nil {
		return nil, err
	}
	t.LastTx = anchor
	t.Owner = w.Signer.Owner()
	if err = w.Signer.SignTransaction(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (w *Wallet) SendTransaction(t *types.Transaction) (*types.Transaction, error) {
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
