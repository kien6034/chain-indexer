package indexer

import (
	"encoding/json"
	"fmt"
)

// Define the structure for the transaction items.
type TxItem struct {
	TxID     string     `json:"txid"`
	Version  int        `json:"version"`
	Locktime int        `json:"locktime"`
	Vin      []VinItem  `json:"vin"`
	Vout     []VoutItem `json:"vout"` // If there are outputs, they would be added here
	Size     int        `json:"size"`
	Weight   int        `json:"weight"`
	Sigops   int        `json:"sigops"`
	Fee      int        `json:"fee"`
	Status   TxStatus   `json:"status"`
}

// Define the structure for previous outputs.
type Prevout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

// Define the structure for outputs (vout), empty in this transaction but included for completeness.
type VoutItem struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

type VinItem struct {
	TxID         string   `json:"txid"`
	Vout         int      `json:"vout"`
	Prevout      Prevout  `json:"prevout"`
	ScriptSig    string   `json:"scriptsig"`
	ScriptSigAsm string   `json:"scriptsig_asm"`
	Witness      []string `json:"witness"`
	IsCoinbase   bool     `json:"is_coinbase"`
	Sequence     uint32   `json:"sequence"`
}

// Define the structure for transaction status.
type TxStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}

// methods for vinItems
func (vin *VinItem) ContainAddress(address string) bool {
	return vin.Prevout.ScriptPubKeyAddress == address
}

// methods for voutItems
func (vout *VoutItem) ContainAddress(address string) bool {
	return vout.ScriptPubKeyAddress == address
}

// methods for tx item
func (tx *TxItem) Print() {
	txItemJSON, err := json.MarshalIndent(tx, "", "    ")
	if err != nil {
		// Handle the error if JSON marshaling fails
		fmt.Println("Error marshaling txItem:", err)
	} else {
		// Print the nicely formatted JSON string
		fmt.Println(string(txItemJSON))
	}
}

func (tx *TxItem) IsConfirmed() bool {
	return tx.Status.Confirmed
}
