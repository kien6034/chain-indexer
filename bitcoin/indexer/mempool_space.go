package indexer

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type MempoolSpaceClient struct {
	isMainnet bool
	baseAPI   string
}

func NewMempoolSpaceClient(isMainnet bool) *MempoolSpaceClient {
	baseAPI := MempoolSpaceAPI + "testnet" + "/api"
	if isMainnet {
		baseAPI = MempoolSpaceAPI + "/api"
	}

	return &MempoolSpaceClient{
		baseAPI:   baseAPI,
		isMainnet: isMainnet,
	}
}

func (c *MempoolSpaceClient) GetAddressTransactions(address string) ([]TxItem, error) {
	url := fmt.Sprintf("%s/address/%s/txs", c.baseAPI, address)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err // Return the error if the request failed
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // Return the error if reading the response body failed
	}

	var txItems []TxItem
	if err := json.Unmarshal(body, &txItems); err != nil {
		return nil, err
	}
	return txItems, nil
}

func (c *MempoolSpaceClient) GetAddressUTXOs(address string) ([]UTXO, error) {
	url := fmt.Sprintf("%s/address/%s/utxo", c.baseAPI, address)
	fmt.Println("url: ", url)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err // Return the error if the request failed
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // Return the error if reading the response body failed
	}

	var utxos []UTXO
	if err := json.Unmarshal(body, &utxos); err != nil {
		return nil, err
	}
	return utxos, nil
}

func (c *MempoolSpaceClient) FetchTransactionScriptPubKey(txid string, vout int, netParams *chaincfg.Params) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/tx/%s/hex", c.baseAPI, txid)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the transaction hex from the response.
	txHex, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the hex to a transaction.
	txBytes, err := hex.DecodeString(string(txHex))
	if err != nil {
		return nil, err
	}
	tx := wire.NewMsgTx(2)
	if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		return nil, err
	}

	// Check if the vout index is within the bounds of the transaction's outputs.
	if vout >= len(tx.TxOut) {
		return nil, fmt.Errorf("vout index out of bounds")
	}

	// Extract the scriptPubKey.
	scriptPubKey := tx.TxOut[vout].PkScript
	return scriptPubKey, nil
}

func (c *MempoolSpaceClient) BroadcastTx(txHex string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tx", c.baseAPI), bytes.NewBufferString(txHex))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
