package dcr

import (
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/chaincfg"
	//"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrjson"
	"github.com/btcsuite/btcd/txscript"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrutil"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrec"

	"github.com/gaozhengxin/cryptocoins/src/go/btc"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var ChainConfig = chaincfg.MainNetParams

var RequiredConfirmations = int64(1)

var allowHighFees = true

var feeRate, _ = dcrutil.NewAmount(0.0001)

var hashType = txscript.SigHashAll

type DCRHandler struct{
	btcHandler *btc.BTCHandler
}

func NewDCRHandler () *DCRHandler {
	return &DCRHandler{
		btcHandler: btc.NewBTCHandlerWithConfig(config.DCR_SERVER_HOST,config.DCR_SERVER_PORT,config.DCR_USER,config.DCR_PASSWD,config.DCR_USESSL),
	}
}

var DCR_DEFAULT_FEE, _ = new(big.Int).SetString("50000",10)

func (h *DCRHandler) GetDefaultFee() *big.Int {
	return DCR_DEFAULT_FEE
}

func (h *DCRHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
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

func (h *DCRHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return h.btcHandler.BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress, amount, jsonstring)
}

func (h *DCRHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	return h.btcHandler.SignTransaction(hash, privateKey)
}

func (h *DCRHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	return h.btcHandler.MakeSignedTransaction(rsv, transaction)
}

func (h *DCRHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	return h.btcHandler.SubmitTransaction(signedTransaction)
}

func (h *DCRHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	return h.btcHandler.GetTransactionInfo(txhash)
}

func (h *DCRHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	return
}

