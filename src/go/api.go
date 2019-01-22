package api

import (
	"math/big"

	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/eos"
	"github.com/gaozhengxin/cryptocoins/src/go/eth"
	"github.com/gaozhengxin/cryptocoins/src/go/erc20"
	"github.com/gaozhengxin/cryptocoins/src/go/xrp"
	"github.com/gaozhengxin/cryptocoins/src/go/trx"
)

type TransactionHandler interface {

	// 公钥to地址
	// eos的address是随机生成的账户名, msg是eos格式的公钥: EOS6JUDHVf4mbrbMNXxhMVJUj5Tz14d1jYpdjC8ZvRgFb4jhrBKEe
	PublicKeyToAddress(pubKeyHex string) (address_or_account_name string, msg string, err error)

	// 构造未签名交易
	// btc, ripple需要fromPublicKey
	// eth, erc20需要fromAddress
	// eos需要账户名(fromAddress)和fromPublicKey
	BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error)

	// 签名函数 txhash 输出 rsv
	//SignTransaction(hash []string, address string) (rsv []string, err error)
	SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error)

	// 构造签名交易
	MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error)

	// 提交交易
	SubmitTransaction(signedTransaction interface{}) (ret string, err error)

	// txhash 查出 fromaddress， toaddress， 交易额
	GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error)

	// 账户查账户余额
	GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error)

	//SetCrypto(cryptoType string)  // ecdsa, ed25519
	//GetCrypto() string
	//SetCurve(curve string)  // k1, r1
	//GetCurve() string
}

func NewTransactionHandler(coinType string) (txHandler TransactionHandler) {
	switch coinType {
	case "BTC":
		return &btc.BTCTransactionHandler{}
	case "EOS":
		return &eos.EOSTransactionHandler{}
	case "ETH":
		return &eth.ETHTransactionHandler{}
	case "ERC20":
		return &erc20.ERC20TransactionHandler{}
	case "XRP":
		return &xrp.XRPTransactionHandler{}
	case "TRX":
		return &trx.TRXTransactionHandler{}
	}
	return nil
}
