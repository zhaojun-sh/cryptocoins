package ltc

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var ChainConfig = chaincfg.Params {
	PubKeyHashAddrID: 0x30,
}

var allowHighFees = true

type LTCHandler struct {
	btcHandler *btc.BTCHandler
}

func NewLTCHandler () *LTCHandler {
	return &LTCHandler{
		btcHandler: btc.NewBTCHandlerWithConfig(config.LTC_SERVER_HOST,config.LTC_SERVER_PORT,config.LTC_USER,config.LTC_PASSWD,config.LTC_USESSL),
	}
}

func (h *LTCHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
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
func (h *LTCHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, jsonstring)
}

// NOT completed, may or not work
func (h *LTCHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error) {
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *LTCHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *LTCHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	c, _ := rpcutils.NewClient(config.LTC_SERVER_HOST,config.LTC_SERVER_PORT,config.LTC_USER,config.LTC_PASSWD,config.LTC_USESSL)
	ret, err= btc.SendRawTransaction (c, signedTransaction.(*btc.AuthoredTx).Tx, allowHighFees)
	return h.SubmitTransaction(signedTransaction)
}

func (h *LTCHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.GetTransactionInfo(txhash)
}

// TODO
func (h *LTCHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	err = fmt.Errorf("function currently not available")
	return nil, err
}
