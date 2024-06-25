package main

import (
	"log"

	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction/data_item"
	"github.com/liteseed/goar/wallet"
)

func SendBundle() {
	w, err := wallet.FromPath("./arweave.json", "https://arweave.net")
	if err != nil {
		log.Fatal(err)
	}

	d := w.CreateDataItem([]byte("test"), "", "", nil)
	_, err = w.SignDataItem(d)
	if err != nil {
		log.Fatal(err)
	}

	b, err := w.CreateBundle(&[]data_item.DataItem{*d})
	if err != nil {
		log.Fatal(err)
	}

	tx := w.CreateTransaction(b.RawData, "", "", &[]tag.Tag{{Name: "test", Value: "test"}})
	log.Println(tx)
	_, err = w.SignTransaction(tx)
	if err != nil {
		log.Fatal(err)
	}
	err = w.SendTransaction(tx)
	if err != nil {
		log.Fatal(err)
	}
}
