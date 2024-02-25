package bitcoin

type BitcoinClient struct {
	*MempoolSpaceClient
}

func NewBitcoinClient(isMainnet bool) *BitcoinClient {
	return &BitcoinClient{
		NewMempoolSpaceClient(isMainnet),
	}
}

func (c *BitcoinClient) GetAddressTransactions(address string) (string, error) {
	return c.MempoolSpaceClient.GetAddressTransactions(address)
}
