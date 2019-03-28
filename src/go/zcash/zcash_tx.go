package zcash

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

type ZCASHTransactionHandler struct {
	btcHandler *btc.BTCTransactionHandler
}

func NewZCASHTransactionHandler () *ZCASHTransactionHandler {
	return *ZCASHTransactionHandler{
		btcHandler = btc.NewBTCHandlerWithConfig(config.ZCASH_SERVER_HOST,config.ZCASH_SERVER_PORT,config.ZCASH_USER,config.ZCASH_PASSWD,config.ZCASH_USESSL)
	}
}

func (h *ZCASHTransactionHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
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
	address = "t" + address
	return
}

// NOT completed, may or not work
func (h *ZCASHTransactionHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	transaction, digests, err = btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, args)
	return
}

// NOT completed, may or not work
func (h *ZCASHTransactionHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *ZCASHTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *ZCASHTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return h.SubmitTransaction(signedTransaction)
}

func (h *ZCASHTransactionHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.GetTransactionInfo(txhash)
}

// TODO
func (h *ZCASHTransactionHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	err = fmt.Errorf("function currently not available")
	return nil, err
}
