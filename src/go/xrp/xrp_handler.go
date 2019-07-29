package xrp

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"github.com/btcsuite/btcd/btcec"
	"github.com/rubblelabs/ripple/crypto"
	"github.com/rubblelabs/ripple/data"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
)

var (
	fee int64 = 1
	url = config.ApiGateways.RippleGateway.ApiAddress
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

type XRPHandler struct {}

func NewXRPHandler () *XRPHandler {
	return &XRPHandler{}
}

var XRP_DEFAULT_FEE, _ = new(big.Int).SetString("1",10)

func (h *XRPHandler) GetDefaultFee() *big.Int {
	return XRP_DEFAULT_FEE
}

func (h *XRPHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
	pub, err := hex.DecodeString(pubKeyHex)
	address = XRP_publicKeyToAddress(pub)
	return
}

// jsonstring:'{"fee":1}'
func (h *XRPHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	var args interface{}
	json.Unmarshal([]byte(jsonstring), &args)
	userFee := args.(map[string]interface{})["fee"]
	if userFee != nil {
		fee = int64(userFee.(float64))
	}
	if fromAddress == "" {
		fromAddress, err = h.PublicKeyToAddress(fromPublicKey)
		if err != nil {
			return
		}
	}
	pub, err := hex.DecodeString(fromPublicKey)
	xrp_pubKey := XRP_importPublicKey(pub)
	amt := amount.String()
	txseq := getSeq(fromAddress)
	transaction, hash, _ := XRP_newUnsignedPaymentTransaction(xrp_pubKey, nil, txseq, toAddress, amt, fee, "", false, false, false)
	digests = append(digests, hash.String())
	return
}

func (h *XRPHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	seed := strings.Split(privateKey.(string), "/")[0]
	keySeqStr := strings.Split(privateKey.(string), "/")[1]
	key := XRP_importKeyFromSeed(seed, "ecdsa")
	ki, err1 := strconv.Atoi(keySeqStr)
	if err1 != nil {
		err = fmt.Errorf("invalid key sequence")
		return
	}
	keyseq := uint32(ki)

	hashBytes, err := hex.DecodeString(hash[0])
	if err != nil {
		return
	}

	sig, err := crypto.Sign(key.Private(&keyseq), hashBytes, nil)
	if err != nil {
		return
	}
	signature, err := btcec.ParseSignature(sig, btcec.S256())
	fmt.Printf("==================\n!!!!! signature is %+v\n==================\n\n", signature)
	if err != nil {
		return
	}
	rx := fmt.Sprintf("%X", signature.R)
	rx = make64(rx)
	sx := fmt.Sprintf("%X", signature.S)
	sx = make64(sx)
	rsv = append(rsv, rx + sx + "00")
	return
}

func (h *XRPHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	sig := rsvToSig(rsv[0])
	signedTransaction = XRP_makeSignedTx(transaction.(data.Transaction), sig)
	return
}

func (h *XRPHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
fmt.Printf("++++++++++++++++++++++++\n%+v\n++++++++++++++++++++++++\n",signedTransaction)
	ret := XRP_submitTx(signedTransaction.(data.Transaction))

	var retStruct interface{}
	json.Unmarshal([]byte(ret), &retStruct)
	result := retStruct.(map[string]interface{})["result"].(map[string]interface{})
	if result["error"] != nil {
		txhash = ""
		err = fmt.Errorf("%v, %v Error message: %v", result["error"], result["error_exception"], result["error_message"])
		return
	}
fmt.Printf("%+v\n\n",result)
	if result["engine_result_message"].(string) == "The transaction was applied. Only final in a validated ledger." {
		txhash = result["tx_json"].(map[string]interface{})["hash"].(string)
	} else if res := result["engine_result_message"].(string); res != "" {
		err = fmt.Errorf(res)
	}

	return
}

func (h *XRPHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	data := "{\"method\":\"tx\", \"params\":[{\"transaction\":\"" + txhash + "\", \"binary\":false}]}"
	ret := rpcutils.DoPostRequest(url, "", data)

	var retStruct interface{}
	json.Unmarshal([]byte(ret), &retStruct)
	result := retStruct.(map[string]interface{})["result"].(map[string]interface{})

	if result["error"] != nil {
		err = fmt.Errorf("%v, error code: %v,  error message: %v", result["error"].(string), result["error_code"].(float64), result["error_message"].(string))
		return
	}

	fromAddress = result["Account"].(string)
	toAddress := result["Destination"].(string)
	amt := result["Amount"].(string)
	transferAmount, _ := new(big.Int).SetString(amt, 10)
	txOutputs = append(txOutputs, types.TxOutput{
		ToAddress: toAddress,
		Amount: transferAmount,
	})
	return
}

func (h *XRPHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	account := getAccount(address)
	balance, _ = new(big.Int).SetString(account.Balance, 10)
	return
}

func parseAccount(s string) *data.Account {
	account, err := data.NewAccountFromAddress(s)
	checkErr(err)
	return account
}

func parseAmount(s string) *data.Amount {
	amount, err := data.NewAmount(s)
	checkErr(err)
	return amount
}

func parsePaths(s string) *data.PathSet {
	ps := data.PathSet{}
	for _, pathStr := range strings.Split(s, ",") {
		path, err := data.NewPath(pathStr)
		checkErr(err)
		ps = append(ps, path)
	}
	return &ps
}

type JsonRet struct {
	Result Account_info_Res
}

func rsvToSig(rsv string) []byte {
	b, _ := hex.DecodeString(rsv)
	rx := hex.EncodeToString(b[:32])
	sx := hex.EncodeToString(b[32:64])
	r, _ := new(big.Int).SetString(rx, 16)
	s, _ := new(big.Int).SetString(sx, 16)
	signature := &btcec.Signature{
		R: r,
		S: s,
	}
	return signature.Serialize()
}

type Account_info_Res struct {
	Account_data Account
}

type Account struct {
	Balance string
	Sequence uint32
}

func make64 (str string) string {
	for l := len(str); l < 64; l++ {
		str = "0" + str
	}
	return str
}

func getAccount (address string) (Account) {
	// TODO
	reader := strings.NewReader("{\"method\":\"account_info\",\"params\":[{\"account\":\"" + address + "\"}]}")
        request, err := http.NewRequest("POST", url, reader)
        checkErr(err)
        client := &http.Client{}
        resp, err := client.Do(request)
        checkErr(err)
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
	jsonRet := new(JsonRet)
	err = json.Unmarshal(body, jsonRet)
	checkErr(err)
	return jsonRet.Result.Account_data
}

// 查帐户目前的sequence
func getSeq(address string) uint32 {
	account := getAccount(address)
	return account.Sequence
}

func XRP_newUnsignedSimplePaymentTransaction(fromAddress string, publicKey []byte, toAddress string, amount *big.Int, fee int64) (data.Transaction, data.Hash256, []byte) {
	dcrm_key := XRP_importPublicKey(publicKey)
	z := new(big.Int).Div(amount, big.NewInt(1000000))
	d := new(big.Int).Sub(amount, new(big.Int).Mul(amount, big.NewInt(1000000)))
	amt := z.String() + "." + d.String() + "/XRP/" + fromAddress
	dcrm_txseq := getSeq(fromAddress)  // 一般是1
	return XRP_newUnsignedPaymentTransaction(dcrm_key, nil, dcrm_txseq, toAddress, amt, fee, "", false, false, false)
}

// 普通xrp转账
func XRP_Remit(seed string, cryptoType string, keyseq *uint32, toaddress string, amount *big.Int, fee int64) {
        key := XRP_importKeyFromSeed(seed, cryptoType)
        fromaddress := XRP_getAddress(key, keyseq)
        txseq := getSeq(fromaddress)
	z := new(big.Int).Div(amount, big.NewInt(1000000))
	d := new(big.Int).Sub(amount, new(big.Int).Mul(z, big.NewInt(1000000)))
	amt := z.String() + "." + d.String() + "/XRP/" + fromaddress
        tx, hash, _ := XRP_newUnsignedPaymentTransaction(key, keyseq, txseq, toaddress, amt, fee, "", false, false, false)
        sig := XRP_getSig(tx, key, keyseq, hash, nil)
        signedTx := XRP_makeSignedTx(tx, sig)
        res := XRP_submitTx(signedTx)
        fmt.Println(res)
}

// 给一个地址打10000块钱激活, 需要一个有足够钱的大帐户
// 大帐户地址: rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE
// 大帐户seed: ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf
// 大帐户密钥类型: ecdsa  keysequence: 0
func XRP_FundAddress(toaddress string) {
        key := XRP_importKeyFromSeed("ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf","ecdsa")
        keyseq := uint32(0)
        txseq := uint32(1)  // 新帐户是1
        tx, hash, msg := XRP_newUnsignedPaymentTransaction(key, &keyseq, txseq, toaddress, "10000/XRP/rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE", int64(10), "", false, false, false)

        // 签名
        sig := XRP_getSig(tx, key, &keyseq, hash, msg)

        // 构造交易结构, 发送交易
        XRP_makeSignedTx(tx, sig)
        res := XRP_submitTx(tx)
        fmt.Printf("%v\n",res)
}

// keyseq is only supported by ecdsa, leave nil when key crypto type is ed25519
// amt format: "value/currency/issuer"
func XRP_newUnsignedPaymentTransaction (key crypto.Key, keyseq *uint32, txseq uint32, dest string, amt string, fee int64, path string, nodirect bool, partial bool, limit bool) (data.Transaction, data.Hash256, []byte) {

	destination, amount := parseAccount(dest), parseAmount(amt)
	payment := &data.Payment{
		Destination: *destination,
		Amount:      *amount,
	}
	payment.TransactionType = data.PAYMENT

	if path != "" {
		payment.Paths = parsePaths(path)
	}
	payment.Flags = new(data.TransactionFlag)
	if nodirect {
		*payment.Flags = *payment.Flags | data.TxNoDirectRipple
	}
	if partial {
		*payment.Flags = *payment.Flags | data.TxPartialPayment
	}
	if limit {
		*payment.Flags = *payment.Flags | data.TxLimitQuality
	}

	base := payment.GetBase()

	base.Sequence = txseq

	fei, err := data.NewNativeValue(fee)
	checkErr(err)
	base.Fee = *fei

	//TODO Set Account
	copy(base.Account[:], key.Id(keyseq))

	payment.InitialiseForSigning()
	copy(payment.GetPublicKey().Bytes(), key.Public(keyseq))
	hash, msg, err := data.SigningHash(payment)
	checkErr(err)

	return payment, hash, msg
}

func XRP_getSig(tx data.Transaction, key crypto.Key, keyseq *uint32, hash data.Hash256, msg []byte) []byte {
	sig, err := crypto.Sign(key.Private(keyseq), hash.Bytes(), append(tx.SigningPrefix().Bytes(), msg...))
	checkErr(err)
	return sig
}

func XRP_makeSignedTx(tx data.Transaction, sig []byte) data.Transaction {
	*tx.GetSignature() = data.VariableLength(sig)
	hash, _, err := data.Raw(tx)
	checkErr(err)
	copy(tx.GetHash().Bytes(), hash.Bytes())
	return tx
}

func XRP_submitTx(signedTx data.Transaction) string {
	_, raw, err := data.Raw(signedTx)
	checkErr(err)
	txBlob := fmt.Sprintf("%X", raw)

	data := "{\"method\":\"submit\",\"params\":[{\"tx_blob\":\"" + txBlob + "\"}]}"

	return rpcutils.DoPostRequest(url, "", data)
}
