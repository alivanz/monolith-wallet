package main

import (
	"encoding/hex"
	"flag"
	"log"
	"strconv"
	"strings"

	crypto "github.com/alivanz/go-crypto"
	"github.com/alivanz/go-crypto/bitcoin"
	microwallet "github.com/alivanz/go-microwallet"
)

const (
	// test
	BitcoinTestnet  = "btc_testnet"
	EthereumGanache = "eth_ganache"

	// Action
	GetAddress = "get-address"
	GetBalance = "get-balance"
	Transfer   = "transfer"
)

var (
	// hack
	privkey string

	coin     string
	rawdests string
	action   string
	windex   int
	wallet   crypto.Wallet

	// transfer destination
	dests       map[string]int64
	dest_change string
)

func main() {
	// Coins
	coins := make(map[string]func())
	coins["btc"] = coin_bitcoin
	coins["btc_testnet"] = coin_bitcoin_testnet
	coins["doge"] = coin_dogecoin
	coins["dast"] = coin_dash
	coins["eth"] = coin_ethereum
	coins["eth_ganache"] = coin_ethereum_ganache

	coins_key := make([]string, 0)
	for coin, _ := range coins {
		coins_key = append(coins_key, coin)
	}

	flag.StringVar(&privkey, "privkey", "", "Override privkey")
	flag.IntVar(&windex, "windex", 0, "Wallet index in micro-wallet")
	flag.StringVar(&coin, "coin", "btc", strings.Join(coins_key, ", "))
	flag.StringVar(&rawdests, "dest", "", "Output destination (format=addr0:value0,addr1:value1,...)")
	flag.StringVar(&action, "action", "", strings.Join([]string{GetAddress, GetBalance, Transfer}, ", "))
	flag.Parse()

	if len(action) == 0 {
		flag.PrintDefaults()
	}

	// Open wallet
	if len(privkey) == 0 {
		bank, err := microwallet.OpenBank(nil)
		if err != nil {
			log.Print(err)
			log.Fatal("Unable to locate micro-wallet")
		}
		wallet, err = bank.Open(windex)
		if err != nil {
			log.Print(err)
			log.Fatal("Unable to open wallet")
		}
	} else {
		data, err := hex.DecodeString(privkey)
		if err != nil {
			log.Print(err)
			log.Fatal("Unable to decode privkey")
		}
		wallet, err = bitcoin.NewWallet(data)
		if err != nil {
			log.Print(err)
			log.Fatal("Unable to open wallet")
		}
	}

	// Parse destination
	dests = make(map[string]int64)
	for _, dest := range strings.Split(rawdests, ",") {
		if len(dest) == 0 {
			continue
		}
		addrvalue := strings.Split(dest, ":")
		if len(addrvalue) != 2 {
			log.Print(dest)
			log.Fatal("wrong output dest format")
		}
		starget := addrvalue[0]
		samount := addrvalue[1]
		if samount == "*" {
			dest_change = starget
			continue
		}
		amount, err := strconv.ParseInt(samount, 10, 64)
		if err != nil {
			log.Print(err)
			log.Fatal("Wrong wallet format")
		}
		dests[starget] = amount
	}

	// Coin decision
	subprogram, found := coins[coin]
	if !found {
		log.Fatal("unknown coin " + coin)
	}
	subprogram()
}
