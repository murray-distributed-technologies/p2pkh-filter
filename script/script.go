package script

import (
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
	pushtx "github.com/murray-distributed-technologies/go-pushtx/script"
)

func AppendFilter(s *bscript.Script) (*bscript.Script, error) {
	var err error
	// strip sha1 hash
	s.AppendOpcodes(bscript.Op1, bscript.OpSPLIT, bscript.OpSWAP, bscript.OpSPLIT, bscript.OpNIP)

	// push x bytes to split after OP_PUSH_TX
	if err = s.AppendPushDataHexString(""); err != nil {
		return nil, err
	}
	s.AppendOpcodes(bscript.OpSPLIT)

	// strip pubkeyHash
	s.AppendOpcodes(bscript.Op2, bscript.OpSPLIT, bscript.Op1, bscript.OpSPLIT, bscript.OpNIP)

	// build template
	s.AppendOpcodes(bscript.OpCAT, bscript.OpCAT)

	//take hash and check
	s.AppendOpcodes(bscript.OpSHA1, bscript.OpEQUALVERIFY)
	return s, nil
}

func NewLockingScript(address string) (*bscript.Script, error) {
	var err error
	s := &bscript.Script{}

	//add hash of script template
	if err = s.AppendPushDataHexString(""); err != nil {
		return nil, err
	}

	// grab preimage
	s.AppendOpcodes(bscript.Op1, bscript.OpPICK)
	// get locking script from preimage
	if s, err = pushtx.AppendGetLockingScriptFromPreimage(s); err != nil {
		return nil, err
	}
	// strip data from template
	if s, err = AppendFilter(s); err != nil {
		return nil, err
	}
	// add pushTX
	if s, err = pushtx.AppendPushTx(s); err != nil {
		return nil, err
	}
	// add p2pkh
	if s, err = pushtx.AppendP2PKH(s, address); err != nil {
		return nil, err
	}
	return s, nil
}

func NewUnlockingScript(pubKey, preimage, sig []byte, sigHashFlag sighash.Flag) (*bscript.Script, error) {
	sigBuf := []byte{}
	sigBuf = append(sigBuf, sig...)
	sigBuf = append(sigBuf, uint8(sigHashFlag))

	scriptBuf := [][]byte{sigBuf, pubKey}

	s := &bscript.Script{}
	err := s.AppendPushDataArray(scriptBuf)
	if err != nil {
		return nil, err
	}

	if preimage != nil {
		if err = s.AppendPushData(preimage); err != nil {
			return nil, err
		}
	}
	return s, nil
}
