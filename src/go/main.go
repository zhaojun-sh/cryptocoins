package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/fusion/go-fusion/crypto"

	"github.com/cryptocoins/src/go/xrp"
)

func main() {
	test_erc20()
	//test_xrp()
}

func test_erc20 () {
	privateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
	pub := crypto.FromECDSAPub(&privateKey.PublicKey)
	pubKeyHex := hex.EncodeToString(pub)
	fmt.Printf("pubKeyHex is %v\n\n", pubKeyHex)

	h := NewTransactionHandler("ERC20")
/*
	address, _, err := h.PublicKeyToAddress(pubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("address is %v\n\n", address)

	address2 := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	transaction, digest, err := h.BuildUnsignedTransaction(address, "", address2, big.NewInt(1), big.NewInt(10), uint64(4000000), "BNB")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n%v\n\n", transaction, digest)

	rsv, err := h.SignTransaction(digest, privateKey)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
                return
	}

	signedTransaction, err := h.MakeSignedTransaction(rsv, transaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n", signedTransaction)

	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%s\n\n", ret)
*/
	fromAddress, toAddress, amount, _, err := h.GetTransactionInfo("0xf9e16303a1b5a59b12e18be82aaed2363621844d8b78961db57d1af7aa89419f")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
                return
	}
	fmt.Printf("from: %v\nto: %v\namount: %v\n\n", fromAddress, toAddress, amount)

	balance, err := h.GetAddressBalance("0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140", "GUSD")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("balance: %v\n\n", balance)
}

func test_xrp () {
	h := NewTransactionHandler("XRP")
	fmt.Printf("h type: %T\n\n", h)

	pubKeyHex := "a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be"
	address, msg, err := h.PublicKeyToAddress(pubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("address is %v\nmsg is \"%s\"\n\n", address, msg)

	//address1 := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"
	address2 := "raF1e6TSKtB34MZ9USrKphQAW5hYbARFWK"
	fromKey := xrp.XRP_importKeyFromSeed("ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf", "ecdsa")
	keyseq := uint32(0)
	fromPubKeyHex := hex.EncodeToString(fromKey.Public(&keyseq))

	fee := int64(10)
	transaction, digest, err := h.BuildUnsignedTransaction("", fromPubKeyHex, address2, big.NewInt(100), &fee)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n%v\n\n", transaction, digest)

	seed := "ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf"
	rsv, err := h.SignTransaction(digest, seed)

	signedTransaction, err := h.MakeSignedTransaction(rsv, transaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n", signedTransaction)

	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%s\n\n", ret)

	fromAddress, toAddress, transferAmount, _, err := h.GetTransactionInfo("738F84EE3BDCA016680916021EF613A2DA3B1302F8D6448E33D1976A55796C54")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}
	fmt.Printf("fromAddress: %v\ntoAddress: %v\ntransferAmount: %v\n\n", fromAddress, toAddress, transferAmount)
}
