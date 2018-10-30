package main

import (
	"log"

	"github.com/alivanz/go-crypto/bitcoin"
)

func coin_bitcoin() {
	compressed := true
	// get pubkey
	pubkey, err := wallet.PubKey()
	if err != nil {
		log.Fatal(err)
	}
	// get address
	addr, err := bitcoin.PubKeyToAddress(bitcoin.MainNetworkDesc.PubKeyHashCode, bitcoin.SerializePubkey(pubkey, compressed))
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case GetAddress:
		log.Print(addr)
	}
}
