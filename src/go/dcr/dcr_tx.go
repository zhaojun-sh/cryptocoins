package dcr

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime/debug"

	"github.com/btcsuite/btcd/btcec"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/chaincfg"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrjson"
	"github.com/btcsuite/btcd/txscript"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrutil"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrec"

	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
)

var ChainConfig = chaincfg.MainNetParams

var RequiredConfirmations = int64(1)

var allowHighFees = true

var feeRate, _ = dcrutil.NewAmount(0.0001)

var hashType = txscript.SigHashAll

type DCRTransactionHandler struct{}

func (h *DCRTransactionHandler) PublicKeyToAddress(pubKeyHex string) (address string, msg string, err error){
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
	pkHash := dcrutil.Hash160(b)
	addressPubKeyHash, err := dcrutil.NewAddressPubKeyHash(pkHash, &ChainConfig, dcrec.STEcdsaSecp256k1)
	if err != nil {
		return
	}
	address = addressPubKeyHash.EncodeAddress()
	return
}

func (h *DCRTransactionHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error) {
	return
}

func (h *DCRTransactionHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return
}

func (h *DCRTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return
}

func (h *DCRTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return
}

func (h *DCRTransactionHandler) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	cmd := dcrjson.NewGetRawTransactionCmd(txhash, nil)

	marshalledJSON, err := dcrjson.MarshalCmd("1.0", 66, cmd)
	if err != nil {
		return
	}

	c, _ := rpcutils.NewClient(config.DCR_SERVER_HOST,config.DCR_SERVER_PORT,config.DCR_USER,config.DCR_PASSWD,config.DCR_USESSL)
	retJSON, err := c.Send(string(marshalledJSON))
	if err != nil {
		return
	}

	var rawTx interface{}
	json.Unmarshal([]byte(retJSON), &rawTx)
	rawTxStr := rawTx.(map[string]interface{})["result"].(string)

	cmd2 := dcrjson.NewDecodeRawTransactionCmd(rawTxStr)

	marshalledJSON2, err := dcrjson.MarshalCmd("1.0", 66, cmd2)
	if err != nil {
		return
	}
	retJSON2, err := c.Send(string(marshalledJSON2))
	if err != nil {
		return
	}
	var tx interface{}
	json.Unmarshal([]byte(retJSON2), &tx)
	toAddress = tx.(map[string]interface{})["result"].(map[string]interface{})["vout"].([]interface{})[0].(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
	flt := tx.(map[string]interface{})["result"].(map[string]interface{})["vout"].([]interface{})[0].(map[string]interface{})["value"].(float64)
	amt, err := dcrutil.NewAmount(flt)
	transferAmount = big.NewInt(int64(amt.ToUnit(dcrutil.AmountAtom)))

	vintx0 := tx.(map[string]interface{})["result"].(map[string]interface{})["vin"].([]interface{})[0].(map[string]interface{})["txid"]
	coinbase := tx.(map[string]interface{})["result"].(map[string]interface{})["vin"].([]interface{})[0].(map[string]interface{})["coinbase"]
	if vintx0 == nil {
		fromAddress = coinbase.(string)
		return
	}
	vintx := vintx0.(string)
	//vintx := tx.(map[string]interface{})["result"].(map[string]interface{})["vin"].([]interface{})[0].(map[string]interface{})["txid"].(string)
	vinvout := int(tx.(map[string]interface{})["result"].(map[string]interface{})["vin"].([]interface{})[0].(map[string]interface{})["vout"].(float64))

	cmd3 := dcrjson.NewGetRawTransactionCmd(vintx, nil)

	marshalledJSON3, err := dcrjson.MarshalCmd("1.0", nil, cmd3)
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

	cmd4 := dcrjson.NewDecodeRawTransactionCmd(rawTxStr2)

	marshalledJSON4, err := dcrjson.MarshalCmd("1.0", nil, cmd4)
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

func (h *DCRTransactionHandler) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error){
	return
}

