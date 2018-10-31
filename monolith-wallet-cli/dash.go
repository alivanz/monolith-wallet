package main

import (
	"log"

	dogecoin "github.com/freddyisman/go-dogecoin"
)

func coin_dash() {
	pubkeyhash := dogecoin.WalletToPubKeyHash(wallet)
	var coin dogecoin.Dash
	coindata := coin.CreateCoin(pubkeyhash)

	switch action {
	case GetAddress:
		log.Print(coindata.Address)
	case GetBalance:
		log.Print(coindata.Balance)
	case Transfer:
		destx := make([]dogecoin.Destination, 0)
		for target, amount := range dests {
			destx = append(destx, dogecoin.Destination{target, uint64(amount)})
		}
		var sendvalue uint64
		for _, outaddr := range destx {
			sendvalue = sendvalue + outaddr.Value
		}
		totalfee, numindex := dogecoin.ChangeUnspent(coindata, sendvalue, &destx)
		if coindata.Balance >= (sendvalue + totalfee) {
			signtx := dogecoin.CreateSignedTransaction(coindata, destx, wallet, numindex)
			log.Printf("signtxhex    : %v\n", signtx)
			// coin.Broadcast(signtx)
		} else if coindata.Balance >= sendvalue {
			log.Printf("total fee belum melewati batas minimum, transaksi tidak dapat dilakukan\n")
		} else {
			log.Printf("saldo tidak mencukupi, transaksi tidak dapat dilakukan\n")
		}
	}

}
