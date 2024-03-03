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
	// create indexer
	indexer := indexer.NewBitcoinClient(false)

	// // Get address transactions
	// txItems, err := indexer.GetAddressTransactions("tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm")

	// if err != nil {
	// 	panic(err)
	// }

	// for _, tx := range txItems {
	// 	tx.VerbalInfo()
	// }

	w := wallet.NewBtcWallet(PRIVATE_KEY, false)

	r, err := w.SendTxWithMemo(*indexer, "tb1qnnuc6efguvx097v74j7udxt05ra90g0txwaar6", 1000, "0x5b5fDd1510F817Ece8bBD911d7028144522c4429", "97")
	if err != nil {
		panic(err)
	}

	log.Printf("response: %s", r)
}
