package woc

import (
	"context"
	"errors"

	"github.com/mrz1836/go-whatsonchain"
)

func GetTransactionOutput(txId string, vout int) (*whatsonchain.VoutInfo, error) {
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)
	txInfo, err := client.GetTxByHash(context.Background(), txId)
	if err != nil {
		return nil, err
	}
	if len(txInfo.Vout) < vout {
		return nil, errors.New("transaction didnt have enough outputs")
	}
	return &txInfo.Vout[vout], nil
}
