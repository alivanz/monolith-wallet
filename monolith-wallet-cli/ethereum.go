package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func coin_ethereum() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	coin_ethereum_(client)
}
func coin_ethereum_ganache() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	coin_ethereum_(client)
}
func coin_ethereum_(client *ethclient.Client) {
	publicKeyECDSA, err := wallet.PubKey()
	if err != nil {
		log.Fatal(err)
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	address := fromAddress.Hex()

	switch action {
	case GetAddress:
		log.Print(address)
	case GetBalance:
		account := common.HexToAddress(address)
		balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			log.Fatal(err)
		}
		fbalance := new(big.Float)
		fbalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
		log.Print(ethValue)
	case Transfer:
		for dest, amount := range dests {
			nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
			if err != nil {
				log.Fatal(err)
			}
			toAddress := common.HexToAddress(dest)
			value := big.NewInt(amount)
			gasLimit := uint64(21000)
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
			chainID, err := client.NetworkID(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			hash := types.NewEIP155Signer(chainID).Hash(tx)

			signature, err := CreateEtherSignature(wallet, hash[:], publicKeyBytes[:])
			if err != nil {
				log.Fatal(err)
			}
			sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Equal(sigPublicKey, publicKeyBytes) != true {
				signstring := hex.EncodeToString(signature[:64])
				signature, err = hex.DecodeString(signstring + "01")
				if err != nil {
					log.Fatal(err)
				}
			}
			// verifikasi signature
			verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signature[:len(signature)-1])
			if verified == true {
				// men-signature transaksi
				signedTx, err := tx.WithSignature(types.NewEIP155Signer(chainID), signature)
				if err != nil {
					log.Fatal(err)
				}
				// mem-broadcast transaksi
				err = client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("tx sent: %s", signedTx.Hash().Hex())
			}
		}
	}
}
