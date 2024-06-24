package main

import (
	"log"

	"github.com/liteseed/goar/wallet"
)

func main() {
	w, err := wallet.FromPath("./arweave.json", "https://arweave.net")
	if err != nil {
		log.Fatal(err)
	}

	tx := w.CreateTransaction(nil, "F7fmxSBJx5RlIRrt825iIEAL110cKP2Bf8tYd0Q1STU", "100", nil)
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
