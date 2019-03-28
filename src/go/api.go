package api

import (
	"math/big"

	"github.com/gaozhengxin/cryptocoins/src/go/types"

	"github.com/gaozhengxin/cryptocoins/src/go/bch"
	"github.com/gaozhengxin/cryptocoins/src/go/bitgold"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/dash"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr"
	"github.com/gaozhengxin/cryptocoins/src/go/eos"
	"github.com/gaozhengxin/cryptocoins/src/go/eth"
	"github.com/gaozhengxin/cryptocoins/src/go/etc"
	"github.com/gaozhengxin/cryptocoins/src/go/erc20"
	"github.com/gaozhengxin/cryptocoins/src/go/ltc"
	"github.com/gaozhengxin/cryptocoins/src/go/trx"
	"github.com/gaozhengxin/cryptocoins/src/go/tether"
	"github.com/gaozhengxin/cryptocoins/src/go/vechain"
	"github.com/gaozhengxin/cryptocoins/src/go/xrp"
	"github.com/gaozhengxin/cryptocoins/src/go/zcash"
)

type TransactionHandler interface {

	// 公钥to dcrm地址
	PublicKeyToAddress(pubKeyHex string) (address string, err error)

	// 构造未签名交易
	BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error)

	// 签名函数 txhash 输出 rsv 测试用
	//SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error)

	// 构造签名交易
	MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error)

	// 提交交易
	SubmitTransaction(signedTransaction interface{}) (txhash string, err error)

	// 根据交易hash查交易信息
	// fromAddress 交易发起方地址
	// txOutputs 交易输出切片, txOutputs[i].ToAddress 第i条交易接收方地址, txOutputs[i].Amount 第i条交易转账金额
	GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error)

	// 账户查账户余额
	GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error)
}

func NewTransactionHandler(coinType string) (txHandler TransactionHandler) {
	switch coinType {
	case "BITGOLD":
		return bitgold.NewBTCTransactionHandler{}
	case "BCH":
		return bch.NewBCHTransactionHandler{}
	case "BTC":
		return btc.NewBTCTransactionHandler{}
	case "DASH":
		return dash.NewDASHTransactionHandler{}
	case "DCR":
		return dcr.NewDCRTransactionHandler{}
	case "EOS":
		return eos.NewEOSTransactionHandler{}
	case "ETH":
		return eth.NewETHTransactionHandler{}
	case "ETC":
		return etc.NewETCTransactionHandler{}
	case "ERC20":
		return erc20.NewERC20TransactionHandler{}
	case "LTC":
		return ltc.NewLTCTransactionHandler{}
	case "TRX":
		return trx.NewTRXTransactionHandler{}
	case "TETHER":
		return tether.NewTETHERTransactionHandler{}
	case "VECHAIN":
		return vechain.NewVECHAINTransactionHandler{}
	case "XRP":
		return xrp.NewXRPTransactionHandler{}
	case "ZCASH":
		return zcash.NewZCASHTransactionHandler{}
	}
	return nil
}
