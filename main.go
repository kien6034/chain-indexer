package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kien6034/chain-indexer/bitcoin/indexer"
	"github.com/kien6034/chain-indexer/bitcoin/wallet"
)

var PRIVATE_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PRIVATE_KEY = os.Getenv("PRIVATE_KEY")
	if PRIVATE_KEY == "" {
		log.Fatal("PRIVATE_KEY environment variable is not set.")
	}
}

func main() {
	isMainnet := false
	wallet := wallet.NewBtcWallet(PRIVATE_KEY, isMainnet)
	client := indexer.NewBitcoinClient(isMainnet) // testnet

	// Get client utxos
	wifAddr, err := wallet.GetWifAddress()
	if err != nil {
		log.Fatal(err)
	}
	client.GetAddressUTXOs(wifAddr)

	// get wallet
}
