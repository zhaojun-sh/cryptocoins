package main

import (
	"math/big"

	"cryptocoins/src/go/erc20"
)

type TransactionHandler interface {

	// 公钥to地址
	PublicKeyToAddress(pubKeyHex string) (address string, err error)

	//GetPublicKey(pubKeyHex string) (publicKey interface{}, err error)

	// 构造未签名交易
	BuildUnsignedTransaction(fromAddress, toAddress string, amount *big.Int, args ...interface{}) (transaction interface{}, digests []string, err error)

	// 构造签名交易
	MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error)

	// 提交交易
	SubmitTransaction(signedTransaction interface{}) (ret string, err error)

	//SetCrypto(cryptoType string)
	//GetCrypto() string
	//SetCurve(curve string)
	//GetCurve() string
	//Canonical() bool
}

func NewTransactionHandler(coinType string) (txHandler TransactionHandler) {
	switch coinType {
	case "ERC20":
		return &erc20.ERC20TransactionHandler{}
/*
	case "BTC":
		return &BTCTransactionHandler{}
	case "ETH":
		return &ETHTransactionHandler{}
	case "ERC20":
		return &ERC20TransactionHandler{}
	case "XRP":
		return &XRPTransactionHandler{}
	case "EOS":
		return &EOSTransactionHandler{}
*/
	}
	return nil
}
