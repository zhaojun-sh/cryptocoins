package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/gaozhengxin/cryptocoins/src/go/eos"
	"github.com/gaozhengxin/cryptocoins/src/go/xrp"
	"github.com/gaozhengxin/cryptocoins/src/go/trx"
	api "github.com/gaozhengxin/cryptocoins/src/go"
)

func main() {
	//test_bch()
	//test_btc()
	//test_dash()
	//test_tether()
	//test_ltc()
	//test_eos()
	//test_eth()
	test_etc()
	//test_vechain()
	//test_erc20()
	//test_xrp()
	//test_tron()
	//test_zcash()
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
/*
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
*/
	fmt.Printf("========== %s ==========\n\n", "test get transaction info")
	fromAddress2, toAddress2, amount, _, err := h.GetTransactionInfo(queryTxHash)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("from: %v\nto: %v\namount: %v\n\n", fromAddress2, toAddress2, amount)
/*
	fmt.Printf("========== %s ==========\n\n", "test get balance")
	balance, err := h.GetAddressBalance(queryAddress, query_balance_args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("balance: %v\n\n", balance)
*/
}

func test_bch () {
	fmt.Printf("=========================\n           BCH           \n=========================\n\n")
	h := api.NewTransactionHandler("BCH")
	fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "f6f24bd236252574c1ba7086ac1178b37abf5041653127bd6b400bae0f0a9b00"
	test_common (h, nil, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
}

func test_dash () {
	fmt.Printf("=========================\n           DASH           \n=========================\n\n")
	h := api.NewTransactionHandler("DASH")
	fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "60a8de0be75d153be34d39d31a4c2b0c6904be1354b462f9af8fe75fc1c2fc5e"
	test_common (h, nil, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
}

func test_tether () {
	fmt.Printf("=========================\n           TETHER           \n=========================\n\n")
	h := api.NewTransactionHandler("TETHER")
	fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "30fdcd207ba3cc3f9f5d69ab25e2290c71c1008564d7c27fdc3d96a797f65643"
	test_common (h, nil, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
}

func test_btc () {
	fmt.Printf("=========================\n           BTC           \n=========================\n\n")
	h := api.NewTransactionHandler("BTC")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	toAddress := "mg1KnRaekxjZbvdUNDKxxJycd3hNbxMomA"
	var build_tx_args []interface{}
	build_tx_args = append(build_tx_args, float64(0), "")
	queryTxHash := "6bf5a5077234908b44f69f5587f92c027a68374d88ccc36012663b4ebcdbc802"
	queryAddress := "2MteNic4ttfvkYCJYEaYMuqrNcnc6xzwoBL"
	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, build_tx_args, queryTxHash, queryAddress, nil)
}

func test_ltc () {
	fmt.Printf("=========================\n           LTC           \n=========================\n\n")
	h := api.NewTransactionHandler("LTC")
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	queryTxHash := "ae1359de01c84c1750faa71ac62ed8381e97fa1156b280861ebb01fc84f538aa"
	test_common (h, nil, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
}

func test_zcash () {
	fmt.Printf("=========================\n           ZCASH           \n=========================\n\n")
	h := api.NewTransactionHandler("ZCASH")
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	queryTxHash := "6a37436e56d4e4ac081c816b628404bc28a216afc0d0514ca0d490cc28fa5a28"
	test_common (h, nil, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
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
	memo := "1234"
	build_tx_args = append(build_tx_args, memo)
	queryTxHash := "0cd1f75fff840bca344d1aa61c4bc5a0d97082b04a8bb8ee4e3a255a86f7cf19"
	queryAcct := "degtjwol11u3"

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAcctName, toAcctName, build_tx_args, queryTxHash, queryAcct, nil)
}

func test_eth () {
	fmt.Printf("=========================\n           ETH           \n=========================\n\n")
	h := api.NewTransactionHandler("ETH")

	//fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
//fromPrivateKey, _ := crypto.HexToECDSA("0ea7b1364bc2d1d58f35324b4c0deaa129cc9bd2728e0942e7592f2836cbb530")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x426B635fD6CdAf5E4e7Bf5B2A2Dd7bc6c7360FBd"

	toAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	var build_tx_args []interface{}
	build_tx_args = append(build_tx_args, big.NewInt(8000000000), uint64(50000))

	queryTxHash := "0xf9e16303a1b5a59b12e18be82aaed2363621844d8b78961db57d1af7aa89419f"

//	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
	queryAddress := fromAddress

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, build_tx_args, queryTxHash, queryAddress, nil)
}

func test_etc () {
	fmt.Printf("=========================\n           ETC           \n=========================\n\n")
	h := api.NewTransactionHandler("ETC")

	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	queryTxHash := "0xb3c033e2f22cd8f7e6f07d3da87bd92cfe8ca6128632d6b7d596f417bc440588"

	test_common (h, fromPrivateKey, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
}

func test_vechain () {
	fmt.Printf("=========================\n           VECHAIN           \n=========================\n\n")
	h := api.NewTransactionHandler("VECHAIN")

	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	queryTxHash := "0x1d4e7f364082341b719aa51cb86c907b37bb2335eca2b4243b6a5de39f87e87c"

	test_common (h, fromPrivateKey, fromPubKeyHex, "", "", nil, queryTxHash, "", nil)
}

func test_erc20 () {
	fmt.Printf("=========================\n           ERC20           \n=========================\n\n")
	h := api.NewTransactionHandler("ERC20")

	//fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
//fromPrivateKey, _ := crypto.HexToECDSA("40d6e64ce085269869b178c23a786e499ff2d6a5334fe45964211d25bea973bf")
fromPrivateKey, _ := crypto.HexToECDSA("0ea7b1364bc2d1d58f35324b4c0deaa129cc9bd2728e0942e7592f2836cbb530")

	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	toAddress := "0x426B635fD6CdAf5E4e7Bf5B2A2Dd7bc6c7360FBd"

	var build_tx_args []interface{}
	build_tx_args = append(build_tx_args, big.NewInt(10000000000), uint64(100000), "BNB")

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
	n = 5
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
