package api

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	rcrypto "github.com/rubblelabs/ripple/crypto"

	httputils "cryptocoins/src/go/trx/httputils"
)

const (
	URL = "https://api.shasta.trongrid.io"
	ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	prefix = byte(0x41)  // MainNet
	//prefix = byte(0xA0)
	TRANSFER_CONTRACT = "TransferContract"
)

type PublicKey struct {
	*ecdsa.PublicKey
}

func PublicKeyToHex (pk *PublicKey) (ret string) {
	b := elliptic.Marshal(crypto.S256(), pk.X, pk.Y)
	ret = hex.EncodeToString(b)
	return
}

func HexToPublicKey(pubKeyHex string) (pk *PublicKey, err error) {
	pub, err := hex.DecodeString(pubKeyHex)
	if len(pub) == 65 {
		x := new(big.Int).SetBytes(pub[1:33])
		y := new(big.Int).SetBytes(pub[33:])
		pk = &PublicKey{&ecdsa.PublicKey{
			Curve: crypto.S256(),
			X: x,
			Y: y,
		}}
	} else if len(pub) == 64 {
		x := new(big.Int).SetBytes(pub[:32])
		y := new(big.Int).SetBytes(pub[32:])
		pk = &PublicKey{&ecdsa.PublicKey{
			Curve: crypto.S256(),
			X: x,
			Y: y,
		}}
	} else {
		err = fmt.Errorf("Invalid public key length %v", len(pub))
	}
	return
}

func (pk *PublicKey) Address() (address, address21 string, err error) {
	data := append(pk.X.Bytes(), pk.Y.Bytes()...)
	sha := crypto.Keccak256(data)
	address21bytes := append([]byte{prefix}, sha[len(sha)-20:]...)
	address21 = hex.EncodeToString(address21bytes)
	address = rcrypto.Base58Encode(address21bytes, ALPHABET)
	return
}

func PublicKeyToAddress (pubKeyHex string) (address, address21 string, err error) {
	pk, err := HexToPublicKey(pubKeyHex)
	if err != nil {
		return
	}
	address, address21, err = pk.Address()
	return
}

type Transaction struct {
	Signature string `json:"signature"`
	TxID string `json:"txID"`
	Raw_data RawData `json:"raw_data"`
	Error string `json:"Error,omitempty"`
}

func (tx *Transaction) MarshalJson() (ret string, err error) {
	b, err := json.Marshal(tx)
	if err == nil {
		ret = string(b)
	}
	return
}

func (tx *Transaction) UnmarshalJson(txjson string) (err error) {
	err = json.Unmarshal([]byte(txjson), tx)
	if err == nil {
		if tx.Error != "" {
			err = fmt.Errorf(tx.Error)
		}
	}
	return
}

type RawData struct {
	Contract []Contract `json:"contract"`
	Ref_block_bytes string `json:"ref_block_bytes"`
	Ref_block_hash string `json:"ref_block_hash"`
	Expiration int64 `json:"expiration"`
	Timestamp int64 `json:"timestamp"`
}

type Contract interface {
}

type Transfer struct {
	Amount *big.Int `json:"amount"`
	Owner_address string `json:"owner_address"`
	To_address string `json:"to_address"`
}

func (tf *Transfer) MarshalJSON() (ret string, err error) {
	b, err := json.Marshal(tf)
	if err == nil {
		ret = string(b)
	}
	return
}

func BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, args ...interface{}) (transaction interface{}, digests []string, err error) {
	if len(fromAddress) != 42 {
		b, err1 := rcrypto.Base58Decode(fromAddress, ALPHABET)
		if err1 != nil {
			err = err1
			return
		}
		fromAddress = hex.EncodeToString(b)
	}

	if len(toAddress) != 42 {
		b, err2 := rcrypto.Base58Decode(toAddress, ALPHABET)
		if err2 != nil {
			err = err2
			return
		}
		toAddress = hex.EncodeToString(b)
	}

	tf := &Transfer{
		Amount: amount,
		Owner_address: fromAddress,
		To_address: toAddress,
	}

	tfJson, err := tf.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}

	ret := httputils.DoCurlRequest(URL, "wallet/createtransaction", tfJson)
	//fmt.Printf("%s\n\n", ret)

	transaction = &Transaction{}
	err = transaction.(*Transaction).UnmarshalJson(ret)

	digest := transaction.(*Transaction).TxID
	digests = append(digests, digest)
	return
}

func SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	hashBytes, err := hex.DecodeString(hash[0])
	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		r, s, err1 := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hashBytes)
		if err1 != nil {
			err = err1
			return
		}
		rx := fmt.Sprintf("%x", r)
		sx := fmt.Sprintf("%x", s)
		if isCanonical(&privateKey.(*ecdsa.PrivateKey).PublicKey, s) {
			//fmt.Printf("canonical: %v\n\n", true)
			//fmt.Printf("len(r) is %v\nlen(s) is %v\n\n", len(rx), len(sx))
			rsv = append(rsv, rx + sx + "00")
			break
		}
		//fmt.Printf("canonical: %v\n\n", false)
		if i == 24 {
			err = fmt.Errorf("couldn't find a canonical signature")
			return
		}
	}
	return
}

func MakeSignedTransaction (rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	signedTransaction = transaction
	signedTransaction.(*Transaction).Signature = rsv[0]
	return
}

func SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	req, err := signedTransaction.(*Transaction).MarshalJson()
	fmt.Printf("request data : %s\n\n", req)
	ret = httputils.DoCurlRequest(URL, "wallet/broadcasttransaction", req)
	var result interface{}
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		panic(err.Error())
	}
	if ok := result.(map[string]interface{})["result"]; ok != nil && ok.(bool) == true {
		ret = fmt.Sprintf("success/%v", signedTransaction.(*Transaction).TxID)
	}
	return
}

func GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, err error) {
	data, err := json.Marshal(struct{
		Value string `json:"value"`
	}{
		Value: txhash,
	})
	reqData := string(data)
	//fmt.Println("reqData ", reqData)
	ret := httputils.DoPostRequest(URL, "walletsolidity/gettransactionbyid", reqData)
	//fmt.Printf("ret %s\n\n", ret)
	tx := &Transaction{}
	tx.UnmarshalJson(ret)

	tf := tx.Raw_data.Contract[0].(map[string]interface{})["parameter"].(map[string]interface{})["value"].(map[string]interface{})

	fromAddress = tf["owner_address"].(string)
	toAddress = tf["to_address"].(string)
	transferAmount = big.NewInt(int64(tf["amount"].(float64)))
	return
}

func GetAddressBalance(address string) (balance *big.Int, err error) {
	// TODO
	return
}

// Canonical signatures are those where 1 <= S <= N/2
func isCanonical(pk *ecdsa.PublicKey, s *big.Int) bool {
	if big.NewInt(1).Cmp(s) != 1 && s.Cmp(new(big.Int).Div(pk.Params().N, big.NewInt(2))) != 1 {
		return true
	}
	return false
}

