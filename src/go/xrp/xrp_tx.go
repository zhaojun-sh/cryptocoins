package xrp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"github.com/rubblelabs/ripple/crypto"
	"github.com/rubblelabs/ripple/data"
)

const (
	url = "https://s.altnet.rippletest.net:51234"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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

type Account_info_Res struct {
	Account_data Account
}

type Account struct {
	Balance string
	Sequence uint32
}

func getAccount (address string) (Account) {
	// TODO
	reader := strings.NewReader("{\"method\":\"account_info\",\"params\":[{\"account\":\"" + address + "\"}]}")
        request, err := http.NewRequest("POST", "https://s.altnet.rippletest.net:51234", reader)
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

func XRP_newUnsignedCrossCurrencyPayment (fromaddress string, pubkey []byte, toaddress string, fromcurrency, tocurrency) (data.Transaction, data.Hash256, []byte) {
	// TODO
	key := XRP_importPublicKey(pubkey)
	path := ""  // getPath(...)
	txseq := getSeq(fromaddress)
	amt := ""
	return XRP_newUnsignedPaymentTransaction(key, nil, txseq, toaddress, amt, int64(10), path, true, true, true)
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

	reader := strings.NewReader("{\"method\":\"submit\",\"params\":[{\"tx_blob\":\"" + txBlob + "\"}]}")
	request, err := http.NewRequest("POST", url, reader)
	checkErr(err)
	client := &http.Client{}
	resp, err := client.Do(request)
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}
