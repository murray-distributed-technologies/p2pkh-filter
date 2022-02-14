package main

import (
	"fmt"
	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/murray-distributed-technologies/p2pkh-filter/transaction"
	"github.com/murray-distributed-technologies/p2pkh-filter/woc"
)

func main() {
	privKey := wif.DecodeWif("")
	changeAddress, _ := bscript.NewAddressFromPublicKey(privKey.PrivKey.PubKey(), true)
	destPrivKey := wif.DecodeWif("")
	address, _ := bscript.NewAddressFromPublicKey(desPrivKey.PrivKey.PubKey(), true)
	var sats uint64
	var vOut uint32

	txId := ""
	vOut = 0
	amount := uint64(5000)

	o, _ := woc.GetTransactionOutput(txId, int(vOut))

	sats = uint64(o.Value * 100000000)
	scriptPubKey, err := bscript.NewFromHexString(o.ScriptPubKey.Hex)
	if err != nil {
		fmt.Println(err)
	}

	txIdBytes, _ := hex.DecodeString(txId)

	utxo := &bt.UTXO{
		TxID:          txIdBytes,
		Vout:          vOut,
		LockingScript: scriptPubKey,
		Satoshis:      sats,
	}

	rawTxString, err := transaction.CreateTransaction(utxo, privKey, address, changeAddress, amount)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rawTxString)

}
