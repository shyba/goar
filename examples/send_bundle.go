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
	log.Println(w.Signer)

	dataItems := []data_item.DataItem{}
	for i := 0; i < 10; i++ {
		d := w.CreateDataItem([]byte("test"), "", "", &[]tag.Tag{{Name: "test", Value: "test"}})
		_, err = w.SignDataItem(d)
		if err != nil {
			log.Fatal(err)
		}
		dataItems = append(dataItems, *d)
	}

	b, err := w.CreateBundle(&dataItems)
	if err != nil {
		log.Fatal(err)
	}

	tx := w.CreateTransaction(b.Raw, "", "", &[]tag.Tag{{Name: "test", Value: "test"}, {Name: "test", Value: "test"}, {Name: "test", Value: "test"}})
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
