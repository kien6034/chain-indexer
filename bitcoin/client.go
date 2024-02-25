package bitcoin

import (
	"fmt"
)

type BitcoinClient struct {
	*MempoolSpaceClient
}

func NewBitcoinClient(isMainnet bool) *BitcoinClient {
	return &BitcoinClient{
		NewMempoolSpaceClient(isMainnet),
	}
}

func (c *BitcoinClient) GetAddressTransactions(address string) ([]Transaction, error) {

	txItems, err := c.MempoolSpaceClient.GetAddressTransactions(address)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	for _, txItem := range txItems {
		// Marshal the txItem into a JSON string with indentation
		tx, err := c.ParseTx(txItem, address)
		if err != nil {
			fmt.Println("Error parsing tx:", err) // handling error
		}

		transactions = append(transactions, *tx)
	}

	return transactions, nil
}

/*
* ParseTxOut parses the transaction and find the transaction that our wallet is sending
 */
func (c *BitcoinClient) ParseTx(txItem TxItem, address string) (*Transaction, error) {
	// TODO: check cache / db. If txItem is already handled, return

	if !txItem.IsConfirmed() {
		return nil, fmt.Errorf("tx is not confirmed")
	}

	// decode vin
	totalSpend := int64(0) // amount that the address spend
	senders := make(map[string]int64)
	for _, vin := range txItem.Vin {
		if vin.ContainAddress(address) {
			totalSpend += vin.Prevout.Value
			continue
		}

		if _, exists := senders[vin.Prevout.ScriptPubKeyAddress]; exists {
			senders[vin.Prevout.ScriptPubKeyAddress] += vin.Prevout.Value
		} else {
			senders[vin.Prevout.ScriptPubKeyAddress] = vin.Prevout.Value
		}
	}

	// Decode vout
	totalReceive := int64(0) // amount that the address receive
	receivers := make(map[string]int64)
	for _, vout := range txItem.Vout {
		if vout.ContainAddress(address) {
			totalReceive += vout.Value // subtract the change, since self sendback
			continue
		}
		// Add or update the receiver's total received amount
		if _, exists := receivers[vout.ScriptPubKeyAddress]; exists {
			receivers[vout.ScriptPubKeyAddress] += vout.Value
		} else {
			receivers[vout.ScriptPubKeyAddress] = vout.Value
		}

		// subtract the sender's total sent amount
		if _, exists := senders[vout.ScriptPubKeyAddress]; exists {
			senders[vout.ScriptPubKeyAddress] -= vout.Value
		}
	}

	var senderList []Sender
	for address, value := range senders {
		senderList = append(senderList, Sender{Address: address, Value: value})
	}

	var receiverList []Receiver
	for address, value := range receivers {
		receiverList = append(receiverList, Receiver{Address: address, Value: value})
	}

	var txType TxType
	// perf: improve this logic
	if totalSpend > 0 {
		txType = OutGoing // address send tx

		// discard other senders
		senderList = []Sender{{Address: address, Value: totalSpend}}
	} else if totalReceive > 0 && totalSpend == 0 {
		txType = Incoming // address receive tx

		// discard other receivers
		receiverList = []Receiver{{Address: address, Value: totalReceive}}
	} else {
		// unknown
		return nil, fmt.Errorf("unknown tx type")
	}

	return &Transaction{
		TxId:         txItem.TxID,
		Sender:       senderList,
		Receivers:    receiverList,
		TotalSpend:   totalSpend,
		TotalReceive: totalReceive,
		TxType:       txType,
	}, nil
}
