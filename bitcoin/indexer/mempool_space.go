package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
