package zec

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
	PubKeyHashAddrID: 0x4b,
}

var allowHighFees = true

type ZECHandler struct {
	btcHandler *btc.BTCHandler
}

func NewZECHandler () *ZECHandler {
	return &ZECHandler{
		btcHandler: btc.NewBTCHandlerWithConfig(config.ZCASH_SERVER_HOST,config.ZCASH_SERVER_PORT,config.ZCASH_USER,config.ZCASH_PASSWD,config.ZCASH_USESSL),
	}
}

var ZEC_DEFAULT_FEE, _ = new(big.Int).SetString("50000",10)

func (h *ZECHandler) GetDefaultFee() *big.Int {
	return ZEC_DEFAULT_FEE
}

func (h *ZECHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
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
func (h *ZECHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	transaction, digests, err = h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, jsonstring)
	return
}

// NOT completed, may or not work
func (h *ZECHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *ZECHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *ZECHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return h.SubmitTransaction(signedTransaction)
}

func (h *ZECHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.GetTransactionInfo(txhash)
}

// TODO
func (h *ZECHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	err = fmt.Errorf("function currently not available")
	return nil, err
}
