package bitgold

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var ChainConfig = chaincfg.Params {
	PubKeyHashAddrID: 0x26,
}

var allowHighFees = true

type BITGOLDHandler struct {
	btcHandler *btc.BTCHandler
}

func NewBITGOLDHandler () *BITGOLDHandler {
	return &BITGOLDHandler{
		btcHandler: btc.NewBTCHandlerWithConfig(config.BITGOLD_SERVER_HOST,config.BITGOLD_SERVER_PORT,config.BITGOLD_USER,config.BITGOLD_PASSWD,config.BITGOLD_USESSL),
	}
}

var BITGOLD_DEFAULT_FEE, _ = new(big.Int).SetString("50000",10)

func (h *BITGOLDHandler) GetDefaultFee() *big.Int {
	return BITGOLD_DEFAULT_FEE
}

func (h *BITGOLDHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error){
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
func (h *BITGOLDHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, jsonstring)
}

// NOT completed, may or not work
func (h *BITGOLDHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *BITGOLDHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *BITGOLDHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return h.btcHandler.SubmitTransaction(signedTransaction)
}

func (h *BITGOLDHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.btcHandler.GetTransactionInfo(txhash)
}

// TODO
func (h *BITGOLDHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error){
	err = fmt.Errorf("function currently not available")
	return nil, err
}

