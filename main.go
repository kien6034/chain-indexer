package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kien6034/chain-indexer/bitcoin/indexer"
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

	// Get address transactions
	txItems, err := indexer.GetAddressTransactions("tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm")

	if err != nil {
		panic(err)
	}

	for _, tx := range txItems {
		tx.VerbalInfo()
	}
}
