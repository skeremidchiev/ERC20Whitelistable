package token

import (
	"sync"
)

// WhitelistInput simple wrapper for WhitelistAddress() inputs
type WhitelistInput struct {
	Address string `json:"address"`
}

type WhitelistMultiInput struct {
	Addresses []WhitelistInput `json:"addresses"`
}

// WhitelistInput simple wrapper for Mint() inputs
type MintInput struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

// TxOutput simple wrapper for WhitelistAddress() and Mint() outputs
type TxOutput struct {
	Address         string `json:"address"`
	TransactionHash string `json:"txHash"`
	OK              bool   `json:"ok"`
}

type TxMultiOutput struct {
	*sync.Mutex
	Transactions []TxOutput `json:"txsHash"`
}

func GetTxMultiOutput() *TxMultiOutput {
	return &TxMultiOutput{
		&sync.Mutex{},
		[]TxOutput{},
	}
}

func (txm *TxMultiOutput) Add(tx *TxOutput) {
	txm.Lock()
	txm.Transactions = append(txm.Transactions, *tx)
	txm.Unlock()
}
