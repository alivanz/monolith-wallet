package main

import "github.com/alivanz/go-crypto/bitcoin"

func coin_litecoin() {
	coin_bitcoinx(bitcoin.NetworkDesc{
		PubKeyHashCode: 0x30,
		ScriptHashCode: 0x32,
	}, "ltc", "main")
}
