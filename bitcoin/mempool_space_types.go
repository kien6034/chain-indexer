package bitcoin

// Define the structure for the transaction items.
type TxItems struct {
	TxID     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int    `json:"locktime"`
	Vin      []struct {
		TxID         string   `json:"txid"`
		Vout         int      `json:"vout"`
		Prevout      prevout  `json:"prevout"`
		ScriptSig    string   `json:"scriptsig"`
		ScriptSigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     uint32   `json:"sequence"`
	} `json:"vin"`
	Vout   []voutItem `json:"vout"` // If there are outputs, they would be added here
	Size   int        `json:"size"`
	Weight int        `json:"weight"`
	Sigops int        `json:"sigops"`
	Fee    int        `json:"fee"`
	Status txStatus   `json:"status"`
}

// Define the structure for previous outputs.
type prevout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

// Define the structure for outputs (vout), empty in this transaction but included for completeness.
type voutItem struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

// Define the structure for transaction status.
type txStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}
