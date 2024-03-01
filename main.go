package main

import (
	"github.com/kien6034/chain-indexer/indexer"
)

func main() {
	client := indexer.NewBitcoinClient(false) // testnet

	// Get the transactions for an address
	// txs, err := client.GetAddressTransactions("bc1qw68npyr7xjr7k7622vnvkus0awjusz4rvvl33l")
	// if err != nil {
	// 	panic(err)
	// }

	// for _, tx := range txs {
	// 	// Print the transaction
	// 	tx.VerbalInfo()
	// }

	// Get client utxos
	client.GetAddressUTXOs("tb1qa75lc8j9ku9jn0mmjd8quakqwycscsxjmlcw0a")
}
