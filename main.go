package main

import (
	"fmt"

	"github.com/kien6034/chain-indexer/bitcoin"
)

func main() {
	client := bitcoin.NewBitcoinClient(false) // testnet

	// Get the transactions for an address
	transactions, err := client.GetAddressTransactions("tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm")
	if err != nil {
		panic(err)
	}
	fmt.Println(transactions)
}
