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
	// Coin
	Bitcoin  = "btc"
	Litecoin = "ltc"
	Dogecoin = "doge"
	Ethereum = "eth"

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

	coin     string = Bitcoin
	rawdests string = ""
	action   string
	windex   int
	wallet   crypto.Wallet

	// transfer destination
	dests       map[string]int64
	dest_change string
)

func main() {
	flag.StringVar(&privkey, "privkey", "", "Override privkey")
	flag.IntVar(&windex, "windex", 0, "Wallet index in micro-wallet")
	flag.StringVar(&coin, "coin", Bitcoin, strings.Join([]string{
		Bitcoin,
		BitcoinTestnet,
		Litecoin,
		Dogecoin,
		Ethereum,
		EthereumGanache,
	}, ", "))
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
	switch coin {
	case Bitcoin:
		coin_bitcoin()
	case BitcoinTestnet:
		coin_bitcoin_testnet()
	case Litecoin:
		coin_litecoin()
	case Dogecoin:
		coin_dogecoin()
	case Ethereum:
		coin_ethereum()
	case EthereumGanache:
		coin_ethereum_ganache()
	default:
		log.Fatal("unknown coin " + coin)
	}
}
