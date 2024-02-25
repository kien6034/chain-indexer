package bitcoin

type Receiver struct {
	Address string `json:"address"`
	Value   int64  `json:"value"`
}

type Sender struct {
	Address string `json:"address"`
	Value   int64  `json:"value"`
}

type Transaction struct {
	Sender       []Sender   `json:"sender"` // TODO: sender could be multiple addresses
	Receivers    []Receiver `json:"receivers"`
	TotalSpend   int64      `json:"send_amount"`
	TotalReceive int64      `json:"receive_amount"`
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
