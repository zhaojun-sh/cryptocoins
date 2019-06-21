package cryptocoins

import (
	"math/big"
	"strings"

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
	"github.com/gaozhengxin/cryptocoins/src/go/omni"
	"github.com/gaozhengxin/cryptocoins/src/go/ven"
	"github.com/gaozhengxin/cryptocoins/src/go/xrp"
	"github.com/gaozhengxin/cryptocoins/src/go/zec"
	"github.com/gaozhengxin/cryptocoins/src/go/atom"
	"github.com/gaozhengxin/cryptocoins/src/go/evt"
)

type CryptocoinHandler interface {

	// 公钥to dcrm地址
	PublicKeyToAddress(pubKeyHex string) (address string, err error)

	// 构造未签名交易
	BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error)

	// 签名函数 txhash 输出 rsv 测试用
	SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error)

	// 构造签名交易
	MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error)

	// 提交交易
	SubmitTransaction(signedTransaction interface{}) (txhash string, err error)

	// 根据交易hash查交易信息../api.go
	// fromAddress 交易发起方地址
	// txOutputs 交易输出切片, txOutputs[i].ToAddress 第i条交易接收方地址, txOutputs[i].Amount 第i条交易转账金额
	GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error)

	// 账户查账户余额
	GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error)

	// 默认交易费用
	GetDefaultFee() *big.Int
}

func NewCryptocoinHandler(coinType string) (txHandler CryptocoinHandler) {
	coinTypeC := strings.ToUpper(coinType)
	switch coinTypeC {
	case "BITGOLD":
		return bitgold.NewBITGOLDHandler()
	case "BCH":
		return bch.NewBCHHandler()
	case "BTC":
		return btc.NewBTCHandler()
	case "DASH":
		return dash.NewDASHHandler()
	case "DCR":
		return dcr.NewDCRHandler()
	case "EOS":
		return eos.NewEOSHandler()
	case "ETH":
		return eth.NewETHHandler()
	case "ETC":
		return etc.NewETCHandler()
	case "LTC":
		return ltc.NewLTCHandler()
	case "TRX":
		return trx.NewTRXHandler()
	case "VEN":
		return ven.NewVENHandler()
	case "XRP":
		return xrp.NewXRPHandler()
	case "ZCASH":
		return zec.NewZECHandler()
	case "ATOM":
		return atom.NewAtomHandler()
	default:
		if isEvt(coinTypeC) {
			return evt.NewEvtHandler(coinTypeC)
		}
		if isErc20(coinTypeC) {
			return erc20.NewERC20TokenHandler(coinTypeC)
		}
		if isOmni(coinType) {
			return omni.NewOMNIPropertyHandler(coinType)
		}
		return nil
	}
}

func isEvt(tokentype string) bool {
	return strings.HasPrefix(tokentype, "EVT")
}
func isErc20(tokentype string) bool {
	return strings.HasPrefix(tokentype,"ERC20")
}

func isOmni(propertyname string) bool {
	return strings.HasPrefix(propertyname,"OMNI")
}
