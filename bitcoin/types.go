package bitcoin

import (
	"encoding/json"
	"fmt"
)

type Receiver struct {
	Address string `json:"address"`
	Value   int64  `json:"value"`
}

type Sender struct {
	Address string `json:"address"`
	Value   int64  `json:"value"`
}

type Transaction struct {
	TxId         string     `json:"txid"`
	Sender       []Sender   `json:"sender"` // TODO: sender could be multiple addresses
	Receivers    []Receiver `json:"receivers"`
	TotalSpend   int64      `json:"total_spend"`
	TotalReceive int64      `json:"total_receive"`
	TxType       TxType     `json:"tx_type"`
}

type TxType int

const (
	Incoming TxType = iota // Incoming transaction
	OutGoing               // Outgoing transaction
)

func (t TxType) String() string {
	switch t {
	case Incoming:
		return "Incoming"
	case OutGoing:
		return "Outgoing"
	default:
		return "Unknown"
	}
}

func (t *Transaction) Print() {
	transactionJSON, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		// Handle the error if JSON marshaling fails
		fmt.Println("Error marshaling transaction:", err)
	} else {
		// Print the nicely formatted JSON string
		fmt.Println(string(transactionJSON))
	}
}

func (t *Transaction) VerbalInfo() {
	fmt.Printf("\n \n --- Tx: %s\n", t.TxId)

	if t.TxType == OutGoing {
		if len(t.Sender) == 0 {
			fmt.Println(" Error: sender is empty")
		}
		sender := t.Sender[0]

		fmt.Printf("->> Outgoing: Sent %d sats from %s\n", sender.Value, sender.Address)
		for _, receiver := range t.Receivers {
			fmt.Printf("Receiver: %s || Amount %d", receiver.Address, receiver.Value)
		}

	}

	if t.TxType == Incoming {
		if len(t.Receivers) == 0 {
			fmt.Println("Error: receiver is empty")
		}
		receiver := t.Receivers[0]

		fmt.Printf("<<- Incoming: Received %d sats to %s\n", receiver.Value, receiver.Address)
		for _, sender := range t.Sender {
			fmt.Printf("Sender: %s || Amount %d", sender.Address, sender.Value)
		}
	}
}
