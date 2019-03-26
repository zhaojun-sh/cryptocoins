package tether

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime/debug"
	"strings"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/tetherjson"
)

var btcHandler = new(btc.BTCTransactionHandler)

var allowHighFees = true

type TETHERTransactionHandler struct {}

func (h *TETHERTransactionHandler) PublicKeyToAddress(pubKeyHex string) (address string, msg string, err error){
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
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(pkHash, &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	address = addressPubKeyHash.EncodeAddress()
	return
}

// NOT completed, may or not work
func (h *TETHERTransactionHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error) {
	transaction, digests, err = btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, args)
	return
}

// NOT completed, may or not work
func (h *TETHERTransactionHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *TETHERTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *TETHERTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	c, _ := rpcutils.NewClient(config.TETHER_SERVER_HOST,config.TETHER_SERVER_PORT,config.TETHER_USER,config.TETHER_PASSWD,config.TETHER_USESSL)
	ret, err= btc.SendRawTransaction (c, signedTransaction.(*btc.AuthoredTx).Tx, allowHighFees)
	return
}

func (h *TETHERTransactionHandler) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	cmd := tetherjson.NewOmniGetTransactionCmd(txhash)
	marshalledJSON, err := btcjson.MarshalCmd(1, cmd)
	if err != nil {
		return
	}
	c, _ := rpcutils.NewClient(config.TETHER_SERVER_HOST,config.TETHER_SERVER_PORT,config.TETHER_USER,config.TETHER_PASSWD,config.TETHER_USESSL)
	retJSON, err := c.Send(string(marshalledJSON))
	if err != nil {
		return
	}

	result, err := parseRPCReturn(retJSON)
	if err != nil {
		return
	}

	senderStr := result.(map[string]interface{})["sendingaddress"]
	recipientStr := result.(map[string]interface{})["referenceaddress"]
	amountStr := result.(map[string]interface{})["amount"]
	propertyid := result.(map[string]interface{})["propertyid"]

	if senderStr != nil {
		fromAddress = senderStr.(string)
	}
	if recipientStr != nil {
		toAddress = recipientStr.(string)
	}
	if amountStr != nil {
		amountStr = strings.Replace(amountStr.(string), ".", "", -1)
		transferAmount, _ = new(big.Int).SetString(amountStr.(string), 10)
	}
	if propertyid != nil {
		if propertyid.(float64) != 1 {
			err = fmt.Errorf("wrong property id: %v", propertyid.(float64))
		}
	}

	return
}

// TODO
func (h *TETHERTransactionHandler) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error){
	err = fmt.Errorf("function currently not available")
	return nil, err
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

