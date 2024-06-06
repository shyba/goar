package goar

import (
	"context"
	"fmt"
	"os"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/types"
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

func (w *Wallet) Owner() string {
	return w.Signer.Owner()
}

func (w *Wallet) SendData(data []byte, tags []types.Tag) (types.Transaction, error) {
	return w.SendDataSpeedUp(data, tags, 0)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendDataSpeedUp(data []byte, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(len(data), nil)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     tags,
		Data:     crypto.Base64Encode(data),
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}

	return w.SendTransaction(tx)
}

// SendTransaction: if send success, should return pending
func (w *Wallet) SendTransaction(transaction *types.Transaction) (types.Transaction, error) {
	uploader, err := w.getUploader(transaction)
	if err != nil {
		return types.Transaction{}, err
	}
	err = uploader.Once()
	return *transaction, err
}

func (w *Wallet) SendTransactionConcurrent(ctx context.Context, concurrentNum int, transaction *types.Transaction) (types.Transaction, error) {
	uploader, err := w.getUploader(transaction)
	if err != nil {
		return types.Transaction{}, err
	}
	err = uploader.ConcurrentOnce(ctx, concurrentNum)
	return *transaction, err
}

func (w *Wallet) getUploader(transaction *types.Transaction) (*client.TransactionUploader, error) {
	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return nil, err
	}
	transaction.LastTx = anchor
	transaction.Owner = w.Owner()
	if err = w.Signer.SignTransaction(transaction); err != nil {
		return nil, err
	}
	return client.CreateUploader(w.Client, transaction, nil)
}
