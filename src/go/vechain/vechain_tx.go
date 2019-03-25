package vechain

import  (
	"crypto/ecdsa"
	//"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"


	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/rpcutils"

)

var (
	err error
)

type VECHAINTransactionHandler struct {
}

func (h *VECHAINTransactionHandler) PublicKeyToAddress (pubKeyHex string) (address string, msg string, err error) {
	data := hexEncPubkey(pubKeyHex[2:])

	pub, err := decodePubkey(data)

	addressStruct := ethcrypto.PubkeyToAddress(*pub)
	address = Address(addressStruct).String()
	return
}

func (h *VECHAINTransactionHandler) BuildUnsignedTransaction (fromAddress, fromPublicKey, toAddress string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error) {
	return
}

func (h *VECHAINTransactionHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	return
}

func (h *VECHAINTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	return
}

func (h *VECHAINTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	return
}

func (h *VECHAINTransactionHandler) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	params := make(map[string][]string)
	params["111"] = append(params["111"], "aaa")
	params["222"] = append(params["222"], "bbb")
	b, err := rpcutils.HttpGet(config.VECHAIN_GATEWAY, "transactions/"+txhash+"/receipt", nil)
	if err != nil {
		return
	}
	var body interface{}
	json.Unmarshal(b, &body)
	transfers0 := body.(map[string]interface{})["outputs"].([]interface{})[0].(map[string]interface{})["transfers"].([]interface{})[0].(map[string]interface{})
	fromAddress = transfers0["sender"].(string)
	toAddress = transfers0["recipient"].(string)
	amt := transfers0["amount"].(string)
	transferAmount, ok := new(big.Int).SetString(amt, 0)
	if !ok {
		err = fmt.Errorf("parse amount value error")
		return
	}
	return
}

func (h *VECHAINTransactionHandler) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error) {
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
