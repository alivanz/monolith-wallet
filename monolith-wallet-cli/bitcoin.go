package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/alivanz/go-crypto/bitcoin"
)

func coin_bitcoin() {
	coin_bitcoin_(bitcoin.MainNetworkDesc)
}
func coin_bitcoin_testnet() {
	coin_bitcoin_(bitcoin.TestnetNetworkDesc)
}
func coin_bitcoin_(networkdesc bitcoin.NetworkDesc) {
	compressed := true
	// get pubkey
	pubkey, err := wallet.PubKey()
	if err != nil {
		log.Fatal(err)
	}
	// get address
	addr, err := bitcoin.PubKeyToAddress(networkdesc.PubKeyHashCode, bitcoin.SerializePubkey(pubkey, compressed))
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case GetAddress:
		log.Print(addr)

	case GetBalance:
		resp, err := BCGetBalance(addr.String())
		if err != nil {
			log.Fatal(err)
		}
		data, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(data))

	case Transfer:
		// Change address
		changeaddr := addr
		if len(dest_change) > 0 {
			changeaddr, err = bitcoin.AddressParseBase58(dest_change)
			if err != nil {
				log.Print(err)
				log.Fatal("unable to decode " + dest_change)
			}
		}
		// Get unspent
		unspent, err := GetUnspent(addr.String())
		if err != nil {
			log.Fatal(err)
		}
		// build
		unused := bitcoin.Satoshi(0)
		builder := bitcoin.NewTXBuilder(2)
		for _, tx := range unspent {
			builder.AddTxIn(tx)
			unused = unused + tx.Value
		}
		for starget, amount := range dests {
			target, err := bitcoin.AddressParseBase58(starget)
			if err != nil {
				log.Print(err)
				log.Fatal("Unable to parse " + starget)
			}
			switch target.AddrType() {
			// Standard P2PKH
			case networkdesc.PubKeyHashCode:
				builder.AddTxOut(bitcoin.NewTxOut(bitcoin.P2PKH(target.PubKeyHash()), bitcoin.Satoshi(amount)))
				unused = unused - bitcoin.Satoshi(amount)
			default:
				log.Fatal("unsupported address type")
			}
		}
		// calc fee
		fee := bitcoin.Satoshi(0)
		if fee == 0 {
			fee = bitcoin.Satoshi(len(builder.RawTransaction())) * 50
			log.Printf("Fee: %v", fee)
		}
		// change
		change := unused - fee
		if change < 0 {
			log.Fatal("insufficient funds")
		} else if change > 0 {
			builder.AddTxOut(bitcoin.NewTxOut(bitcoin.P2PKH(changeaddr.PubKeyHash()), bitcoin.Satoshi(change)))
		}
		// raw tx
		log.Print("Raw transaction " + hex.EncodeToString(builder.RawTransaction()))
		hash1 := sha256.Sum256(builder.RawTransaction())
		hash2 := sha256.Sum256(hash1[:])
		// msg hash
		hashes := builder.MessageHashes()
		log.Printf("List msg hash %v", hashes)
		// signature
		scripts := make([][]byte, len(hashes))
		for i, hash := range hashes {
			r, s, err := wallet.Sign(hash)
			if err != nil {
				log.Fatal(err)
			}
			script := bitcoin.P2PKHScriptSig(r, s, bitcoin.SerializePubkey(pubkey, compressed))
			scripts[i] = script
		}
		log.Printf("List scripts %v", scripts)
		// final transaction
		signed, err := builder.SignedTransaction(scripts)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Signed transaction " + hex.EncodeToString(signed))
		hash1 = sha256.Sum256(signed)
		hash2 = sha256.Sum256(hash1[:])
		log.Printf("txhash %s", hex.EncodeToString(hash2[:]))
	}
}
