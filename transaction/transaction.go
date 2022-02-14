package transaction

import (
	"bytes"
	"context"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/sighash"
	btunlocker "github.com/libsv/go-bt/v2/unlocker"
	"github.com/murray-distributed-technologes/p2pkh-filter/script"
	pushtxpreimage "github.com/murray-distributed-technologies/go-pushtx/preimage"
)

func CreateTransaction(utxo *bt.UTXO, privKey *bec.PrivateKey, address, changeAddress string, satoshis uint64) (string, error) {
	var err error
	tx := bt.NewTx()

	if err = tx.FromUTXOS(utxo); err != nil {
		return "", err
	}
	if tx, err = AddOutput(tx, address, satoshis); err != nil {
		return "", err
	}

	if utxo.LockingScript.IsP2PKH() {
		fq := bt.NewFeeQuote()
		if err = tx.ChangeToAddress(changeAddress, fq); err != nil {
			return "", err
		}
	}
	if !utxo.LockingScript.IsP2PKH() {
		lockingScript, err := bscript.NewP2PKHFromAddress(changeAddress)
		if err != nil {
			return "", err
		}
		// need to fix this in go-bt library... just estimating fee for now
		amount := (utxo.Satoshis - satoshis - 500)
		changeOutput := bt.Output{
			Satoshis:      amount,
			LockingScript: lockingScript,
		}
		tx.AddOutput(&changeOutput)
	}

	unlocker := Getter{PrivateKey: privKey}

	// sign input

	if err = tx.FillAllInputs(context.Background(), &unlocker); err != nil {
		return "", err
	}
	return tx.String(), nil

}

func AddOutput(tx *bt.Tx, address string, satoshis uint64) (*bt.Tx, error) {
	lockingScript, err := script.NewLockingScript(address)
	if err != nil {
		return nil, err
	}

	output := bt.Output{
		Satoshis:      satoshis,
		LockingScript: lockingScript,
	}
	tx.AddOutput(&output)
	return tx, nil
}

type Getter struct {
	PrivateKey *bec.PrivateKey
}

func (g *Getter) Unlocker(ctx context.Context, lockingScript *bscript.Script) (bt.Unlocker, error) {
	if lockingScript.IsP2PKH() {
		return &btunlocker.Simple{PrivateKey: g.PrivateKey}, nil
	}
	return &UnlockTx{PrivateKey: g.PrivateKey}, nil
}

type UnlockTx struct {
	PrivateKey *bec.PrivateKey
}

func (u *UnlockTx) UnlockingScript(ctx context.Context, tx *bt.Tx, params bt.UnlockerParams) (*bscript.Script, error) {
	if params.SigHashFlags == 0 {
		params.SigHashFlags = sighash.AllForkID
	}

	// use low s value for preimage to use optimized preimage script
	preimage, err := tx.CalcInputPreimage(params.InputIdx, params.SigHashFlags)
	if err != nil {
		return nil, err
	}
	preimage, nLockTime, err := pushtxpreimage.CheckForLowS(preimage)
	if err != nil {
		return nil, err
	}
	tx.LockTime = nLockTime

	var defaultHex = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var sh []byte
	sh = crypto.Sha256d(preimage)

	if bytes.Equal(defaultHex, preimage) {
		sh = preimage
	}
	sig, err := u.PrivateKey.Sign(sh)
	if err != nil {
		return nil, err
	}

	pubKey := u.PrivateKey.PubKey().SerialiseCompressed()
	signature := sig.Serialise()

	uscript, err := script.NewUnlockingScript(pubKey, preimage, signature, payams.SigHashFlags)
	if err != nil {
		return nil, err
	}

	return uscript, nil
}
