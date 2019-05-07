package bch

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	addrconv "github.com/schancel/cashaddr-converter/address"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var allowHighFees = true

var chainconfig = &chaincfg.TestNet3Params

type BCHHandler struct {
	btcHandler *btc.BTCHandler
}

func NewBCHHandler () *BCHHandler {
	return &BCHHandler{
		btcHandler: btc.NewBTCHandlerWithConfig(config.BCH_SERVER_HOST,config.BCH_SERVER_PORT,config.BCH_USER,config.BCH_PASSWD,config.BCH_USESSL),
	}
}

var BCH_DEFAULT_FEE, _ = new(big.Int).SetString("50000",10)

func (h *BCHHandler) GetDefaultFee() *big.Int {
	return BCH_DEFAULT_FEE
}

func (h *BCHHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error){
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
	legaddr := addressPubKeyHash.EncodeAddress()  // legacy format
	addr, err := addrconv.NewFromString(legaddr)
	if err != nil {
		return
	}
	cashAddress, _ := addr.CashAddress()  // bitcoin cash
	address, err = cashAddress.Encode()
// for lockin test
//************
	address = "qq25efya3nwtmmplsj6j7whzmkj70z609v0rdeq5zf"
//************
	return
}

// NOT completed, may or not work
func (h *BCHHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	transaction, digests, err = h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, jsonstring)
	return
}

// NOT completed, may or not work
func (h *BCHHandler) SignTransaction(hash []string, wif interface{}) (rsv []string, err error){
	return h.btcHandler.SignTransaction(hash, wif)
}

// NOT completed, may or not work
func (h *BCHHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error){
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

// NOT completed, may or not work
func (h *BCHHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return h.btcHandler.SubmitTransaction(signedTransaction)
}

func (h *BCHHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.btcHandler.GetTransactionInfo(txhash)
}

// TODO
func (h *BCHHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error){
	err = fmt.Errorf("function currently not available")
	return nil, err
}
