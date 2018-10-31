package main

import (
	"encoding/hex"
	"time"

	"github.com/alivanz/go-crypto/bitcoin"
)

type BCAddress struct {
	Address            string `json:"address"`
	TotalReceived      int    `json:"total_received"`
	TotalSent          int    `json:"total_sent"`
	Balance            int    `json:"balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	FinalBalance       int    `json:"final_balance"`
	NTx                int    `json:"n_tx"`
	UnconfirmedNTx     int    `json:"unconfirmed_n_tx"`
	FinalNTx           int    `json:"final_n_tx"`
}

type BCAddressExt struct {
	Address            string          `json:"address"`
	TotalReceived      bitcoin.Satoshi `json:"total_received"`
	TotalSent          bitcoin.Satoshi `json:"total_sent"`
	Balance            bitcoin.Satoshi `json:"balance"`
	UnconfirmedBalance bitcoin.Satoshi `json:"unconfirmed_balance"`
	FinalBalance       bitcoin.Satoshi `json:"final_balance"`
	NTx                int             `json:"n_tx"`
	UnconfirmedNTx     int             `json:"unconfirmed_n_tx"`
	FinalNTx           int             `json:"final_n_tx"`
	Txrefs             []BCTxref       `json:"txrefs"`
	TxURL              string          `json:"tx_url"`
}

type BCTxref struct {
	TxHash_       string          `json:"tx_hash"`
	BlockHeight   int             `json:"block_height"`
	TxInputN      int             `json:"tx_input_n"`
	TxOutputN     int             `json:"tx_output_n"`
	Value         bitcoin.Satoshi `json:"value"`
	RefBalance    bitcoin.Satoshi `json:"ref_balance"`
	Spent         bool            `json:"spent"`
	Confirmations int             `json:"confirmations"`
	Confirmed     time.Time       `json:"confirmed"`
	DoubleSpend   bool            `json:"double_spend"`
	Script_       string          `json:"script"`
}

func (tx BCTxref) TxHash() string {
	return tx.TxHash_
}
func (tx BCTxref) OutputIndex() int {
	return tx.TxOutputN
}
func (tx BCTxref) Script() []byte {
	if script, err := hex.DecodeString(tx.Script_); err != nil {
		panic(err)
	} else {
		return script
	}
}

func BCGetBalance(network, address string) (*BCAddress, error) {
	var resp BCAddress
	err := HTTPGet("https://api.blockcypher.com/v1/btc/"+network+"/addrs/"+address+"/balance", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func GetUnspent(network, address string) ([]BCTxref, error) {
	var resp BCAddressExt
	err := HTTPGet("https://api.blockcypher.com/v1/btc/"+network+"/addrs/"+address+"?unspentOnly=true&includeScript=true", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Txrefs, nil
}
