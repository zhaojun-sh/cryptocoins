package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cryptocoins/src/go/eos"
	"github.com/cryptocoins/src/go/xrp"
	"github.com/cryptocoins/src/go/trx"
	api "github.com/cryptocoins/src/go"
)

func main() {
	test_btc()
	//test_eos()
	//test_eth()
	//test_erc20()
	//test_xrp()
	//test_tron()
}

func test_common (h api.TransactionHandler, fromPrivateKey interface{}, fromPubKeyHex, fromAddress, toAddress string, build_tx_args []interface{}, queryTxHash, queryAddress string, query_balance_args []interface{}) {
	fmt.Printf("========== %s ==========\n\n", "test pubkey to address/account_name")
	address, msg, err := h.PublicKeyToAddress(fromPubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("address from pubKeyHex is %v\n\n", address)
	if msg != "" {
		fmt.Printf("msg is %s\n\n", msg)
	}

	fmt.Printf("========== %s ==========\n\n", "test build unsigned transfer transaction")
	transaction, digest, err := h.BuildUnsignedTransaction(fromAddress, fromPubKeyHex, toAddress, big.NewInt(1), build_tx_args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("transaction: %+v\n\ndigest: %v\n\n", transaction, digest)

	fmt.Printf("========== %s ==========\n\n", "test sign with private key")
	rsv, err := h.SignTransaction(digest, fromPrivateKey)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("rsv is %+v\n\n", rsv)

	fmt.Printf("========== %s ==========\n\n", "test make signed transaction")
	signedTransaction, err := h.MakeSignedTransaction(rsv, transaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%+v\n\n", signedTransaction)

	fmt.Printf("========== %s ==========\n\n", "test submit transaction")
	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%s\n\n", ret)

	fmt.Printf("========== %s ==========\n\n", "test get transaction info")
	fromAddress, toAddress, amount, _, err := h.GetTransactionInfo(queryTxHash)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("from: %v\nto: %v\namount: %v\n\n", fromAddress, toAddress, amount)

	fmt.Printf("========== %s ==========\n\n", "test get balance")
	balance, err := h.GetAddressBalance(queryAddress, query_balance_args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("balance: %v\n\n", balance)

}

func test_btc () {
	fmt.Printf("=========================\n           BTC           \n=========================\n\n")
	h := api.NewTransactionHandler("BTC")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	//fromPubKeyHex := "03c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fd"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	toAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	var build_tx_args []interface{}
	build_tx_args = append(build_tx_args, float64(0), "")
	queryTxHash := "26fdbccea16ffec1bdc3a5becb44cda518c1edc0633b94c0cb9339285bf19c9c"
	queryAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, build_tx_args, queryTxHash, queryAddress, nil)
}

func test_eos () {
	fmt.Printf("=========================\n           EOS           \n=========================\n\n")
	h := api.NewTransactionHandler("EOS")
	fromPrivateKey := "5JqBVZS4shWHBhcht6bn3ecWDoZXPk3TRSVpsLriQz5J3BKZtqH"
	fromPubKey := "EOS7EXiYEgNaxgc8ABX5YTATs4fC9nEuCa9fna61X2nZ8Z8KDEMLg"
	fromPubKeyHex, _ := eos.PubKeyToHex(fromPubKey)
	fromAcctName := "gzx123454321"
	toAcctName := "degtjwol11u3"
	var build_tx_args []interface{}
	memo := "hi there"
	build_tx_args = append(build_tx_args, memo)
	queryTxHash := "d8205bbb6e4c131c8a2bc3c6092f8416a291a453b597ddfbc8bae75060678331"
	queryAcct := "degtjwol11u3"

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAcctName, toAcctName, build_tx_args, queryTxHash, queryAcct, nil)
}

func test_eth () {
	fmt.Printf("=========================\n           ETH           \n=========================\n\n")
	h := api.NewTransactionHandler("ETH")

	fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	toAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	var build_tx_args []interface{}
	build_tx_args = append(build_tx_args, big.NewInt(1), uint64(4000000))

	queryTxHash := "0xf9e16303a1b5a59b12e18be82aaed2363621844d8b78961db57d1af7aa89419f"

//	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
	queryAddress := fromAddress

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, build_tx_args, queryTxHash, queryAddress, nil)
}

func test_erc20 () {
	fmt.Printf("=========================\n           ERC20           \n=========================\n\n")
	h := api.NewTransactionHandler("ERC20")

	fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	toAddress := "0xA8dC61209400C9A23bf1fe625c2919c3626Bc157"

	var build_tx_args []interface{}
	build_tx_args = append(build_tx_args, big.NewInt(1), uint64(4000000), "BNB")

	queryTxHash := "0xf9e16303a1b5a59b12e18be82aaed2363621844d8b78961db57d1af7aa89419f"

	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
	var query_balance_args []interface{}
	query_balance_args = append(query_balance_args, "GUSD")

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, build_tx_args, queryTxHash, queryAddress, query_balance_args)
}

func test_xrp () {
	fmt.Printf("=========================\n           XRP           \n=========================\n\n")
	h := api.NewTransactionHandler("XRP")
	fromKey := xrp.XRP_importKeyFromSeed("ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf", "ecdsa")
	keyseq := uint32(0)
	fromPubKeyHex := hex.EncodeToString(fromKey.Public(&keyseq))
	fromAddress := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"
	toAddress := "raF1e6TSKtB34MZ9USrKphQAW5hYbARFWK"
	var build_tx_args []interface{}
	fee := int64(10)
	build_tx_args = append(build_tx_args, &fee)
	fromPrivateKey := "ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf/0"
	queryTxHash := "738F84EE3BDCA016680916021EF613A2DA3B1302F8D6448E33D1976A55796C54"
	queryAddress := "raF1e6TSKtB34MZ9USrKphQAW5hYbARFWK"

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, build_tx_args, queryTxHash, queryAddress, nil)
}

type Seed struct {}

func (s *Seed) Read(p []byte) (n int, err error) {
	n = 1
	return
}

func test_tron() {
	fmt.Printf("=========================\n           TRX           \n=========================\n\n")
	h := api.NewTransactionHandler("TRX")
	fromPrivateKey, _ := ecdsa.GenerateKey(crypto.S256(), &Seed{})
	fromPubKeyHex := trx.PublicKeyToHex(&trx.PublicKey{&fromPrivateKey.PublicKey})
	fromAddress := "417e5f4552091a69125d5dfcb7b8c2659029395bdf"
	toAddress := "41062ae7be408a0cd83a1cb44874d1e748e374d50c"
	queryTxHash := "4ec77591400e30c729f10a11d6b33ef75ea565d46b6314f2fa8d170a4b8f74e1"
	queryAddress := fromAddress

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, nil, queryTxHash, queryAddress, nil)
}
