package bitcoin

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type MempoolSpaceClient struct {
	isMainnet bool
	baseAPI   string
}

func NewMempoolSpaceClient(API string, isMainnet bool) *MempoolSpaceClient {

	baseAPI := MempoolSpaceAPI + "testnet"
	if isMainnet {
		baseAPI = MempoolSpaceAPI + "mainet"
	}

	return &MempoolSpaceClient{
		baseAPI:   baseAPI,
		isMainnet: isMainnet,
	}
}

func (c *MempoolSpaceClient) GetAddressTransactions(address string) (string, error) {
	url := fmt.Sprintf("%s/address/%s", c.baseAPI, address)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err // Return the error if the request failed
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err // Return the error if reading the response body failed
	}

	return string(body), nil
}
