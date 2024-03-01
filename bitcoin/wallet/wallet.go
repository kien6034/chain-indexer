package wallet

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type BtcWallet struct {
	pk       string
	chainCfg chaincfg.Params
}

func NewBtcWallet(pk string, isMainnet bool) *BtcWallet {
	cfg := chaincfg.TestNet3Params
	if isMainnet {
		cfg = chaincfg.MainNetParams
	}

	return &BtcWallet{pk: pk, chainCfg: cfg}
}

func (w *BtcWallet) CreateTx(destination string, amount int64) (string, error) {
	rawTx, err := CreateTx(w.pk, destination, amount)
	if err != nil {
		return "", err
	}
	return rawTx, nil
}

func (w *BtcWallet) GetWifAddress() (string, error) {
	wif, err := btcutil.DecodeWIF(w.pk)

	if err != nil {
		return "", err
	}

	addrWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed()), &w.chainCfg)
	if err != nil {
		return "", err
	}

	return addrWitnessPubKeyHash.EncodeAddress(), nil
}
