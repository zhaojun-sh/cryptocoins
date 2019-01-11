package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type Seed struct {}

func (s *Seed) Read(p []byte) (n int, err error) {
	n = 1
	return
}

func main () {
	// 用固定的seed生成公私钥对
	privKey, err := ecdsa.GenerateKey(crypto.S256(), &Seed{})
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return
	}

	pubKeyHex := publicKeyToHex(&PublicKey{&privKey.PublicKey})
	fmt.Printf("pubKeyHex is %v\nlen is %v\n\n", pubKeyHex, len(pubKeyHex))
	// pubKeyHex = "04E657EE43DDADBA5AF1F1AD4E8098D996C2F3C397E807C9255B0850EA2D151D050F4A0D6451231DD35F0FF653F166C35BCAA0E817520B4DB87DE7E060A72D578E"
	address, addressHex, err := PublicKeyToAddress(pubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}
	fmt.Printf("%v\n%v\n\n", address, addressHex)
	// TMVQGm1qAQYVdetCeGRRkTWYYrLXuHK2HC
	// 417e5f4552091a69125d5dfcb7b8c2659029395bdf

	toAddressHex := "41062ae7be408a0cd83a1cb44874d1e748e374d50c"

	// 构建无签名交易
	tx, digests, err := BuildUnsignedTransaction(addressHex, "", toAddressHex, big.NewInt(10))
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n%v\n\n", tx, digests)

	// 签名
	rsv, err := SignTransaction(digests, privKey)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%v\n\n", rsv)

	// 构建签名交易
	tx, err = MakeSignedTransaction(rsv, tx)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}

	ret, err := SubmitTransaction(tx)
	fmt.Printf("return: %s\n\n", ret)

	// 查看交易
	from, to, amt, err := GetTransactionInfo("646a6c6a7e60dc614dcf2bc35234b62b758e52cdba0381e45e3125f5715bbbf4")
	//from, to, amt, err := GetTransactionInfo(tx.(*Transaction).TxID)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("fromaddress is %v\ntoaddress is %v\namuont is %v\n\n", from, to, amt)
}
