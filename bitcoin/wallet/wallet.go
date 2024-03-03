package wallet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"sort"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/kien6034/chain-indexer/bitcoin/indexer"
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

func (w *BtcWallet) SendTxWithMemo(c indexer.BitcoinClient, receiverAddress string, amount int64, evmAddress string, chainId string) (string, error) {
	memo := fmt.Sprintf("%s:%s", evmAddress, chainId)

	// get wif pubkey address
	p2wpkhAddress, err := w.GetWifPubkeyAddress()
	if err != nil {
		return "", err
	}

	// get utxos
	utxos, err := c.GetAddressUTXOs(p2wpkhAddress)
	if err != nil {
		return "", err
	}

	for _, utxo := range utxos {
		fmt.Printf("UTXO: %s, Vout: %d, Value: %d\n", utxo.TxID, utxo.Vout, utxo.Value)
	}

	fmt.Printf("Using UTXO %s\n: %d", utxos[0].TxID, utxos[0].Vout)

	// define tx input
	tx, err := w.GreedyCoinSelection(utxos, amount)
	if err != nil {
		return "", err
	}

	// get receiver pk script
	receiverPkScript, err := getPkScript(receiverAddress, &w.chainCfg)
	if err != nil {
		return "", err
	}

	// get spender pk script
	spenderPkScript, err := getPkScript(p2wpkhAddress, &w.chainCfg)
	if err != nil {
		return "", err
	}

	// add memo (OP_RETURN) to the tx
	b := txscript.NewScriptBuilder()
	b.AddOp(txscript.OP_RETURN)
	b.AddData([]byte(memo))

	memoScript, err := b.Script()
	if err != nil {
		return "", err
	}

	// add tx output with op_return
	tx.AddTxOut(wire.NewTxOut(0, memoScript))
	// add tx output to transfer sats to receiver address
	tx.AddTxOut(wire.NewTxOut(amount, receiverPkScript))
	// add tx output to transfer unspent sats to spender address
	tx.AddTxOut(wire.NewTxOut(int64(utxos[0].Value)-amount-indexer.RelayFee, spenderPkScript))

	// Fetch utxo[0] script
	script, err := c.FetchTransactionScriptPubKey(utxos[0].TxID, utxos[0].Vout, &w.chainCfg)
	if err != nil {
		return "", err
	}

	log.Printf("\nScript: %s\n", hex.EncodeToString(script))

	wif, err := btcutil.DecodeWIF(w.pk)
	if err != nil {
		return "", err
	}

	// define fetcher (required for witnessSignature)
	fetcher := txscript.NewMultiPrevOutFetcher(nil)

	for _, txIn := range tx.TxIn {
		fetcher.AddPrevOut(txIn.PreviousOutPoint, &wire.TxOut{
			Value:    utxos[0].Value,
			PkScript: script,
		})
	}

	// sign tx
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	witnessSignature, err := txscript.WitnessSignature(tx, sigHashes, 0, int64(utxos[0].Value), script, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		return "", nil
	}

	// since there is only one input, and want to add
	// signature to it use 0 as index
	tx.TxIn[0].Witness = witnessSignature

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)

	// broadcast tx
	res, err := c.BroadcastTx(hex.EncodeToString(signedTx.Bytes()))
	if err != nil {
		return "", err
	}

	return res, nil
}

func (w *BtcWallet) GreedyCoinSelection(utxos []indexer.UTXO, amount int64) (*wire.MsgTx, error) {
	totalAmount := int64(0)
	selectedUtxos := []indexer.UTXO{}

	// Sort UTXOs by value in descending order
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value > utxos[j].Value
	})

	for _, utxo := range utxos {
		if totalAmount >= amount {
			break
		}
		selectedUtxos = append(selectedUtxos, utxo)
		totalAmount += utxo.Value
	}

	if totalAmount < amount {
		return nil, fmt.Errorf("insufficient funds: have %d, need %d", totalAmount, amount)
	}

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs for each selected UTXO
	for _, utxo := range selectedUtxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, fmt.Errorf("invalid txid %s: %v", utxo.TxID, err)
		}
		outPoint := wire.NewOutPoint(hash, uint32(utxo.Vout))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	return tx, nil
}

func getPkScript(address string, chainCfg *chaincfg.Params) ([]byte, error) {
	addr, err := btcutil.DecodeAddress(address, chainCfg)
	if err != nil {
		return nil, err
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return pkScript, nil
}
