package dash

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime/debug"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var ChainConfig = chaincfg.Params {
	PubKeyHashAddrID: 0x4b,
}

var allowHighFees = true

type DASHTransactionHandler struct {
	btcHandler *btc.BTCTransactionHandler
}

func NewDASHTransactionHandler () *DASHTransactionHandler {
	return *DASHTransactionHandler{
		btcHandler = btc.NewBTCHandlerWithConfig(config.DASH_SERVER_HOST,config.DASH_SERVER_PORT,config.DASH_USER,config.DASH_PASSWD,config.DASH_USESSL)
	}
}

func (h *DASHTransactionHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error){
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
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(pkHash, &ChainConfig)
	if err != nil {
		return
	}
	address = addressPubKeyHash.EncodeAddress()
	return
}

// NOT completed, may or not work
func (h *DASHTransactionHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, args)
}

// NOT completed, may or not work
func (h *DASHTransactionHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *DASHTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *DASHTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return h.btcHandler.SubmitTransaction(signedTransaction)
}

func (h *DASHTransactionHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.btcHandler.GetTransactionInfo(txhash)
}

// TODO
func (h *DASHTransactionHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	err = fmt.Errorf("function currently not available")
	return nil, err
}
