package ven

import  (
	"crypto/ecdsa"
	//"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"


	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var (
	err error
)

type VENHandler struct {
}

func NewVENHandler () *VENHandler {
	return &VENHandler{}
}

var VEN_DEFAULT_FEE, _ = new(big.Int).SetString("50000",10)

func (h *VENHandler) GetDefaultFee() *big.Int {
	return VEN_DEFAULT_FEE
}

func (h *VENHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
	data := hexEncPubkey(pubKeyHex[2:])

	pub, err := decodePubkey(data)

	addressStruct := ethcrypto.PubkeyToAddress(*pub)
	address = Address(addressStruct).String()
	return
}

func (h *VENHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return
}

func (h *VENHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	return
}

func (h *VENHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	return
}

func (h *VENHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	return
}

func (h *VENHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	b, err := rpcutils.HttpGet(config.VECHAIN_GATEWAY, "transactions/"+txhash+"/receipt", nil)
	if err != nil {
		return
	}
	var body interface{}
	json.Unmarshal(b, &body)
	transfers := body.(map[string]interface{})["outputs"].([]interface{})[0].(map[string]interface{})["transfers"].([]interface{})
	for _, transfer := range transfers {
		fromAddress = transfer.(map[string]interface{})["sender"].(string)
		toAddress := transfer.(map[string]interface{})["recipient"].(string)
		amt := transfer.(map[string]interface{})["amount"].(string)
		transferAmount, ok := new(big.Int).SetString(amt, 0)
		txOutputs = append(txOutputs, types.TxOutput{
			ToAddress: toAddress,
			Amount: transferAmount,
		})
		if !ok {
			err = fmt.Errorf("parse amount value error")
			return
		}
	}
	return
}

func (h *VENHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	return
}

func hexEncPubkey(h string) (ret [64]byte) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	if len(b) != len(ret) {
		panic("invalid length")
	}
	copy(ret[:], b)
	return ret
}

func decodePubkey(e [64]byte) (*ecdsa.PublicKey, error) {
	p := &ecdsa.PublicKey{Curve: ethcrypto.S256(), X: new(big.Int), Y: new(big.Int)}
	half := len(e) / 2
	p.X.SetBytes(e[:half])
	p.Y.SetBytes(e[half:])
	if !p.Curve.IsOnCurve(p.X, p.Y) {
		return nil, errors.New("invalid secp256k1 curve point")
	}
	return p, nil
}
