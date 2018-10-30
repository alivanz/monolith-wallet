package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	crypto "github.com/alivanz/go-crypto"
)

// fungsi untuk membuat signature
func CreateEtherSignature(wallet crypto.Wallet, hash, publicKeyBytes []byte) ([]byte, error) {
	var signature bytes.Buffer
	r, s, _ := wallet.Sign(hash)
	r_hex := fmt.Sprintf("%x", r)
	s_hex := fmt.Sprintf("%x", s)
	for len(r_hex) < 64 {
		r_hex = "0" + r_hex
	}
	for len(s_hex) < 64 {
		s_hex = "0" + s_hex
	}
	signature.WriteString(r_hex)
	signature.WriteString(s_hex)
	signature.WriteString("00")
	return hex.DecodeString(signature.String())
}
