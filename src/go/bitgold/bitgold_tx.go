package bitgold

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
	PubKeyHashAddrID: 0x26,
}

var allowHighFees = true

type BITGOLDTransactionHandler struct {
	btcHandler *btc.BTCTransactionHandler
}

func NewBITGOLDTransactionHandler () *BITGOLDTransactionHandler {
	return *BITGOLDTransactionHandler{
		btcHandler = btc.NewBTCHandlerWithConfig(config.BITGOLD_SERVER_HOST,config.BITGOLD_SERVER_PORT,config.BITGOLD_USER,config.BITGOLD_PASSWD,config.BITGOLD_USESSL)
	}
}

func (h *BITGOLDTransactionHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error){
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
func (h *BITGOLDTransactionHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, args)
}

// NOT completed, may or not work
func (h *BITGOLDTransactionHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *BITGOLDTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *BITGOLDTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return h.btcHandler.SubmitTransaction(signedTransaction)
}

func (h *BITGOLDTransactionHandler) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	return h.btcHandler.GetTransactionInfo(txhash)
}

// TODO
func (h *BITGOLDTransactionHandler) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error){
	err = fmt.Errorf("function currently not available")
	return nil, err
}

