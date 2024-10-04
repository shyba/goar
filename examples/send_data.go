package main

import (
	"log"

	"github.com/liteseed/goar/wallet"
)

func SendData() {
	w, err := wallet.FromPath("./arweave.json", "http://localhost:1984")
	if err != nil {
		log.Fatal(err)
	}

	tx := w.CreateTransaction([]byte("test"), "", "", nil)
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
