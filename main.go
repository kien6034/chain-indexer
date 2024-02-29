package main

import (
	"github.com/kien6034/chain-indexer/bitcoin"
)

func main() {
	client := bitcoin.NewBitcoinClient(true) // testnet

	// Get the transactions for an address
	txs, err := client.GetAddressTransactions("bc1qw68npyr7xjr7k7622vnvkus0awjusz4rvvl33l")
	if err != nil {
		panic(err)
	}

	for _, tx := range txs {
		// Print the transaction
		tx.VerbalInfo()
	}
}
