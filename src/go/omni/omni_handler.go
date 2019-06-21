package omni

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime/debug"
	"strings"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var allowHighFees = true

var chainconfig = &chaincfg.TestNet3Params

var Properties map[string]string = map[string]string{
	"OMNIOmni":"1",
	"OMNITest Omni":"2",
	"OMNITetherUS":"100",  // TetherUS id on testnet
}

var (
	omnihost = config.ApiGateways.OmniGateway.Host
	omniport = config.ApiGateways.OmniGateway.Port
	omniuser = config.ApiGateways.OmniGateway.User
	omnipasswd = config.ApiGateways.OmniGateway.Passwd
	omniusessl = config.ApiGateways.OmniGateway.Usessl
)

type OmniHandler struct {
	propertyName string
	btcHandler *btc.BTCHandler
}

func NewOMNIHandler () *OmniHandler {
	return &OmniHandler{
		btcHandler: btc.NewBTCHandlerWithConfig(config.ApiGateways.OmniGateway.Host,config.ApiGateways.OmniGateway.Port,config.ApiGateways.OmniGateway.User,config.ApiGateways.OmniGateway.Passwd,config.ApiGateways.OmniGateway.Usessl),
	}
}

func NewOMNIPropertyHandler (propertyname string) *OmniHandler {
	if Properties[propertyname] == "" {
		return nil
	}
	return &OmniHandler{
		propertyName: propertyname,
		btcHandler: btc.NewBTCHandlerWithConfig(config.ApiGateways.OmniGateway.Host,config.ApiGateways.OmniGateway.Port,config.ApiGateways.OmniGateway.User,config.ApiGateways.OmniGateway.Passwd,config.ApiGateways.OmniGateway.Usessl),
	}
}

var OMNI_DEFAULT_FEE, _ = new(big.Int).SetString("10",10)

func (h *OmniHandler) GetDefaultFee() *big.Int {
	return OMNI_DEFAULT_FEE
}

func (h *OmniHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
	if pubKeyHex[:2] == "0x" || pubKeyHex[:2] == "0X" {
		pubKeyHex = pubKeyHex[2:]
	}
	bb, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return
	}
	pubKey, err := btcec.ParsePubKey(bb, btcec.S256())
	if err != nil {
		return
	}
	b := pubKey.SerializeCompressed()
	pkHash := btcutil.Hash160(b)
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(pkHash, chainconfig)
	if err != nil {
		return
	}
	address = addressPubKeyHash.EncodeAddress()
	return
}

// Not Supported
func (h *OmniHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return
}

// NOT completed, may or not work
func (h *OmniHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *OmniHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *OmniHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	c, _ := rpcutils.NewClient(omnihost, omniport, omniuser, omnipasswd, omniusessl)
	ret, err= btc.SendRawTransaction (c, signedTransaction.(*btc.AuthoredTx).Tx, allowHighFees)
	return
}

func (h *OmniHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()

	client, _ := rpcutils.NewClient(omnihost, omniport, omniuser, omnipasswd, omniusessl)
	reqstr := `{"jsonrpc":"1.0","id":"1","method":"omni_gettransaction","params":["`+txhash+`"]}`
	ret, err1 := client.Send(reqstr)
	if err1 != nil {
		err = err1
		return
	}

	if ret == "" {
		err = fmt.Errorf("failed get transaction")
		return
	}
	fmt.Println("GetTransaction: "+ret)
	omniTx := DecodeOmniTx(ret)
	if omniTx.Error != nil {
		err = omniTx.Error
		return
	}

	if h.propertyName != omniTx.PropertyName {
		err = fmt.Errorf("")
	}

	fromAddress = omniTx.From
	txOutput := types.TxOutput{
		ToAddress:omniTx.To,
		Amount:omniTx.Amount,
	}
	txOutputs = append(txOutputs, txOutput)
	return
}

func (h *OmniHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	propertyId := Properties[h.propertyName]
	client, _ := rpcutils.NewClient(omnihost, omniport, omniuser, omnipasswd, omniusessl)
	reqstr := `{"jsonrpc":"1.0","id":"1","method":"omni_getbalance","params":["`+address+`",`+propertyId+`]}`

	ret, err1 := client.Send(reqstr)
	if err1 != nil {
		err = err1
		return
	}
	fmt.Println("GetBalance: "+ret)
	var retObj interface{}
	json.Unmarshal([]byte(ret), &retObj)

	result := retObj.(map[string]interface{})["result"]
	balanceStr := result.(map[string]interface{})["balance"].(string)
	reservedStr := result.(map[string]interface{})["reserved"].(string)
	frozenStr := result.(map[string]interface{})["frozen"].(string)
	balanceStr = strings.Replace(balanceStr,".","",-1)
	reservedStr = strings.Replace(reservedStr,".","",-1)
	frozenStr = strings.Replace(frozenStr,".","",-1)

	balance, _ = new(big.Int).SetString(balanceStr,10)
	reserved, _ := new(big.Int).SetString(reservedStr,10)
	frozen, _ := new(big.Int).SetString(frozenStr,10)

	balance.Sub(balance.Sub(balance,reserved),frozen)

	return
}

func parseRPCReturn (retJSON string) (result interface{}, err error) {
	var ret interface{}
	json.Unmarshal([]byte(retJSON), &ret)
	result = ret.(map[string]interface{})["result"]
	if result == nil {
		errStr := ret.(map[string]interface{})["error"]
		if errStr == nil {
			err = fmt.Errorf("unknown error")
			return
		}
		err = fmt.Errorf(errStr.(string))
	}
	return
}

