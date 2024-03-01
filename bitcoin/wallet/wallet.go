package wallet

type BtcWallet struct {
	pk string
}

func NewBtcWallet(pk string) *BtcWallet {
	return &BtcWallet{pk: pk}
}

func (w *BtcWallet) CreateTx(destination string, amount int64) (string, error) {
	rawTx, err := CreateTx(w.pk, destination, amount)
	if err != nil {
		return "", err
	}
	return rawTx, nil
}
