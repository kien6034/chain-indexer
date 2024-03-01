package wallet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
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

// func (w *BtcWallet) CreateTx(destination string, amount int64) (string, error) {
// 	rawTx, err := CreateTx(w.pk, destination, amount)
// 	if err != nil {
// 		return "", err
// 	}
// 	return rawTx, nil
// }

func (w *BtcWallet) GetWifPubkeyAddress() (string, error) {
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

func (w *BtcWallet) Bench32ToPubkeyScript(bech32Addr string) (string, error) {
	// Decode the Bech32 address
	hrp, decoded, err := bech32.Decode(bech32Addr)
	if err != nil {
		log.Fatalf("Error decoding Bech32 address: %v", err)
	}

	// Convert from Bech32 encoding to a witness program
	converted, err := bech32.ConvertBits(decoded, 5, 8, false)
	if err != nil {
		log.Fatalf("Error converting Bech32 decoded data: %v", err)
	}

	// Ensure the decoded address is for the correct network (testnet in this case)
	if hrp != "tb" {
		log.Fatalf("Address is not a testnet address")
	}

	// Create a P2PKH script
	pkScript, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).AddData(converted).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).Script()
	if err != nil {
		log.Fatalf("Error creating P2PKH script: %v", err)
	}

	// convert back to hexadecimal format
	hexScript := hex.EncodeToString(pkScript)
	fmt.Println("P2PKH Script:", hexScript)

	return hexScript, nil
}

func (w *BtcWallet) SignTx(pkScript string, redeemTx *wire.MsgTx) (string, error) {
	wif, err := btcutil.DecodeWIF(w.pk)
	if err != nil {
		return "", err
	}

	sourcePKScript, err := hex.DecodeString(pkScript)
	if err != nil {
		return "", nil
	}

	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index
	signature, err := txscript.SignatureScript(redeemTx, 0, sourcePKScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return "", nil
	}

	// since there is only one input, and want to add
	// signature to it use 0 as index
	redeemTx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
