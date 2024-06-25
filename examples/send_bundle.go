package main

import (
	"log"

	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction/data_item"
	"github.com/liteseed/goar/wallet"
)

func main() {
	w, err := wallet.FromPath("./arweave.json", "http://localhost:1984")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(w.Signer)

	d1 := w.CreateDataItem([]byte("test"), "", "", nil)
	_, err = w.SignDataItem(d1)
	if err != nil {
		log.Fatal(err)
	}

	b, err := w.CreateBundle(&[]data_item.DataItem{*d1})
	if err != nil {
		log.Fatal(err)
	}

	d2 := w.CreateDataItem(b.Raw, "", "", nil)
	_, err = w.SignDataItem(d2)
	if err != nil {
		log.Fatal(err)
	}

	tx := w.CreateTransaction(d2.Raw, "", "", &[]tag.Tag{{Name: "test", Value: "test"}})
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
