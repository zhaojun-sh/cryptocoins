package ltc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
)

var allowHighFees = true

type LTCTransactionHandler struct{
	btchandler *btc.BTCTransactionHandler
}

func (h *LTCTransactionHandler) PublicKeyToAddress(pubKeyHex string) (address string, msg string, err error){
	return h.btchandler.PublicKeyToAddress(pubKeyHex)
}

// NOT completed, may or not work
func (h *LTCTransactionHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error) {
	return h.btchandler.BuildUnsignedHandler(fromAddress, fromPublicKey, toAddress, amount, args)
}

// NOT completed, may or not work
func (h *LTCTransactionHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btchandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *LTCTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btchandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *LTCTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	btc.SendRawTransaction (signedTransaction.(*btc.AuthoredTx).Tx, allowHighFees)
}

func (h *LTCTransaction) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	cmd := btcjson.NewGetRawTransactionCmd(txhash, nil)

	marshalledJSON, err := btcjson.MarshalCmd(1, cmd)
	if err != nil {
		return
	}

	c, _ := rpcutils.NewClient(config.LTC_SERVER_HOST,config.LTC_SERVER_PORT,config.LTC_USER,config.config.LTC_PASSWD,config.LTC_USESSL)
	retJSON, err := c.Send(string(marshalledJSON))
	if err != nil {
		return
	}

	var rawTx interface{}
	json.Unmarshal([]byte(retJSON), &rawTx)
	rawTxStr := rawTx.(map[string]interface{})["result"].(string)

	cmd2 := btcjson.NewDecodeRawTransactionCmd(rawTxStr)

	marshalledJSON2, err := btcjson.MarshalCmd(1, cmd2)
	if err != nil {
		return
	}
	retJSON2, err := c.Send(string(marshalledJSON2))
fmt.Printf("%v\n\n", retJSON2)
	var tx interface{}
	json.Unmarshal([]byte(retJSON2), &tx)
	toAddress = tx.(map[string]interface{})["result"].(map[string]interface{})["vout"].([]interface{})[0].(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
	flt := tx.(map[string]interface{})["result"].(map[string]interface{})["vout"].([]interface{})[0].(map[string]interface{})["value"].(float64)
	amt, err := btcutil.NewAmount(flt)
	transferAmount = big.NewInt(int64(amt.ToUnit(btcutil.AmountSatoshi)))

	vintx := tx.(map[string]interface{})["result"].(map[string]interface{})["vin"].([]interface{})[0].(map[string]interface{})["txid"].(string)
	vinvout := int(tx.(map[string]interface{})["result"].(map[string]interface{})["vin"].([]interface{})[0].(map[string]interface{})["vout"].(float64))

	cmd3 := btcjson.NewGetRawTransactionCmd(vintx, nil)

	marshalledJSON3, err := btcjson.MarshalCmd(1, cmd3)
	if err != nil {
		return
	}

	retJSON3, err := c.Send(string(marshalledJSON3))
	if err != nil {
		return
	}

	var rawTx2 interface{}
	json.Unmarshal([]byte(retJSON3), &rawTx2)
	rawTxStr2 := rawTx2.(map[string]interface{})["result"].(string)

	cmd4 := btcjson.NewDecodeRawTransactionCmd(rawTxStr2)

	marshalledJSON4, err := btcjson.MarshalCmd(1, cmd4)
	if err != nil {
		return
	}

	retJSON4, err := c.Send(string(marshalledJSON4))
	if err != nil {
		return
	}

	var tx2 interface{}
	json.Unmarshal([]byte(retJSON4), &tx2)

	fromAddress = tx2.(map[string]interface{})["result"].(map[string]interface{})["vout"].([]interface{})[vinvout].(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)

	return
}

// TODO
func (h *LTCTransaction) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error){
	err := fmt.Errorf("function currently not available")
	return nil, err
}

