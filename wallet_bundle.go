package goar

import (
	"context"
	"errors"

	"github.com/liteseed/goar/tx"
)

func (w *Wallet) SendBundleTxSpeedUp(ctx context.Context, concurrentNum int, bundleBinary interface{}, tags []tx.Tag, txSpeed int64) (tx.Transaction, error) {
	bundleTags := []tx.Tag{
		{Name: "Bundle-Format", Value: "binary"},
		{Name: "Bundle-Version", Value: "2.0.0"},
	}
	// check tags cannot include bundleTags Name
	mmap := map[string]struct{}{
		"Bundle-Format":  {},
		"Bundle-Version": {},
	}
	for _, tag := range tags {
		if _, ok := mmap[tag.Name]; ok {
			return tx.Transaction{}, errors.New("tags can not set bundleTags")
		}
	}

	txTags := make([]tx.Tag, 0)
	txTags = append(bundleTags, tags...)
	return w.SendDataConcurrentSpeedUp(ctx, concurrentNum, bundleBinary, txTags, txSpeed)
}

func (w *Wallet) SendBundleTx(ctx context.Context, concurrentNum int, bundleBinary []byte, tags []tx.Tag) (tx.Transaction, error) {
	return w.SendBundleTxSpeedUp(ctx, concurrentNum, bundleBinary, tags, 0)
}
