package bitcoin

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
		baseAPI = MempoolSpaceAPI + "mainnet" + "/api"
	}

	return &MempoolSpaceClient{
		baseAPI:   baseAPI,
		isMainnet: isMainnet,
	}
}

func (c *MempoolSpaceClient) GetAddressTransactions(address string) (string, error) {
	url := fmt.Sprintf("%s/address/%s/txs", c.baseAPI, address)

	fmt.Println("url", url)
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err // Return the error if the request failed
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err // Return the error if reading the response body failed
	}

	var txItems []TxItems
	if err := json.Unmarshal(body, &txItems); err != nil {
		return "", err
	}
	return fmt.Sprintf("%+v", txItems), nil
}
