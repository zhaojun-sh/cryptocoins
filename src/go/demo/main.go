package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/binance-chain/go-sdk/keys"

	"github.com/gaozhengxin/cryptocoins/src/go/eos"
	"github.com/gaozhengxin/cryptocoins/src/go/xrp"
	"github.com/gaozhengxin/cryptocoins/src/go/trx"
	api "github.com/gaozhengxin/cryptocoins/src/go"
)

/*func main () {
	h := api.NewCryptocoinHandler("ERC20GUSD")
	fmt.Printf("========== %s ==========\n\n", "test get transaction info")
	fromAddress, txOutputs, jsonstring, err := h.GetTransactionInfo("0x8b5992bf34ab4c30665631d69e84f6a861ec60375d9a8752fa91bed27b435075")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("from: %v\ntxOutputs: %v\njsonstring: %v\n", fromAddress, txOutputs, jsonstring)
}*/


func main() {
	//test_bch()
	//test_bitgold()
	test_bnb()
	//test_btc()
	//test_dash()
	//test_dcr()
	//test_omni()
	//test_ltc()
	//test_eos()
	//test_eth()
	//test_etc()
	//test_vechain()
	//test_erc20()
	//test_xrp()
	//test_tron()
	//test_zcash()
	//test_atom()
	//test_evt()
	//fmt.Println("Starting main function")
	//select{}
}

var xx = big.NewInt(1)

func test_common (h api.CryptocoinHandler, fromPrivateKey interface{}, fromPubKeyHex, fromAddress, toAddress string, value int64, build_tx_args string, queryTxHash, queryAddress string, query_balance_args string) {
///*
	fmt.Printf("========== %s ==========\n\n", "test pubkey to address/account_name")
	address, err := h.PublicKeyToAddress(fromPubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("address from pubKeyHex is %v\n\n", address)
//*/
/*
	fmt.Printf("========== %s ==========\n\n", "test build unsigned transfer transaction")
	transaction, digest, err := h.BuildUnsignedTransaction(fromAddress, fromPubKeyHex, toAddress, big.NewInt(value), build_tx_args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("transaction: %+v\n\ndigest: %v\n\n", transaction, digest)
*/
/*
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
*/
/*
	fmt.Printf("========== %s ==========\n\n", "test submit transaction")
	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%s\n\n", ret)

*/
///*
	fmt.Printf("========== %s ==========\n\n", "test get transaction info")
	fromAddress2, txOutputs, jsonstring, err := h.GetTransactionInfo(queryTxHash)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("from: %v\ntxOutputs: %v\njsonstring: %v\n", fromAddress2, txOutputs, jsonstring)
//*/
///*
	fmt.Printf("========== %s ==========\n\n", "test get balance")
	balance, err := h.GetAddressBalance(queryAddress, query_balance_args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("balance: %v\n\n", balance)
//*/
}

func test_bitgold () {
	fmt.Printf("=========================\n           BITGOLD           \n=========================\n\n")
	h := api.NewCryptocoinHandler("BITGOLD")
	fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "6e4765956be5b6fe4c4b43f156a07a9662ca591b2a88b09e89d32b51681eb00b"
	test_common (h, nil, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_bch () {
	fmt.Printf("=========================\n           BCH           \n=========================\n\n")
	h := api.NewCryptocoinHandler("BCH")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	//fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "f6f24bd236252574c1ba7086ac1178b37abf5041653127bd6b400bae0f0a9b00"
	test_common (h, fromPrivateKey, fromPubKeyHex, "", "", 1, "", queryTxHash, "qq25efya3nwtmmplsj6j7whzmkj70z609v0rdeq5zf", "")
}

func test_bnb () {
	fmt.Printf("=========================\n           BNB           \n=========================\n\n")
	h := api.NewCryptocoinHandler("BNB")
	km, _ := keys.NewMnemonicKeyManager("govern cancel early excite other fox canvas satoshi social shiver version inch correct web soap always water wine grid fashion voyage finish canal subject")
	fromPrivateKey := km.GetPrivKey()
	fromPubKeyHex := "046377500e127f76b8403aa4e15346dbcde31c54c3914939c8216271ff159fbab5b0713615dbcf7afd262a27f4f65cbf5c5c0e41e13749fcb2956a3a008fcfc939"
	fromAddress := "tbnb1sgudr8w8kfjz0lyuf44sr0x6y29wc96hgm8r5d"
	toAddress := "tbnb1sgudr8w8kfjz0lyuf44sr0x6y29wc96hgm8r5d"
	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 1, "", "B8D69EF425E200C08EBEAB72331E6D6E0D6695B32803F4EB1418FC780F2BFAB4", "tbnb1sgudr8w8kfjz0lyuf44sr0x6y29wc96hgm8r5d", "")
}

func test_dash () {
	fmt.Printf("=========================\n           DASH           \n=========================\n\n")
	h := api.NewCryptocoinHandler("DASH")
	fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "60a8de0be75d153be34d39d31a4c2b0c6904be1354b462f9af8fe75fc1c2fc5e"
	test_common (h, nil, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_dcr () {
	fmt.Printf("=========================\n           DCR           \n=========================\n\n")
	h := api.NewCryptocoinHandler("DCR")
	fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "561541db309e70c7399c6028a8249a09c1da0cee517162de34dca660023c9c6c"
	test_common (h, nil, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_omni () {
	fmt.Printf("=========================\n           OMNI           \n=========================\n\n")
	//h := api.NewCryptocoinHandler("OMNITetherUS") //OMNITetherUS
	h := api.NewCryptocoinHandler("OMNIOmni")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	toAddress := "mt2ciKvkbMcmeWeGWMm1mohxStKpn8SmCY"
	//fromPubKeyHex := "032f7d0667c2f0989dfb588dedc70edfbc5aefdc02304b10a2c58105f8fe3ce38c"
	queryTxHash := "5f2d4a16a2976f0147d73a602deec24abd81af0b15dd91131a258cbb6b12aa11"
	queryAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 1000, "", queryTxHash, queryAddress, "")
}

func test_btc () {
	fmt.Printf("=========================\n           BTC           \n=========================\n\n")
	h := api.NewCryptocoinHandler("BTC")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	toAddress := "mi7z7Nb9H2eWf6Do2qEZn68bAYavjm8fJ2"
//toAddress := "moneyqMan7uh8FqdCA2BV5yZ8qVrc9ikLP"
	build_tx_args := `{"feeRate":0.0001}`
	queryTxHash := "ec12e399ddc6740ef4565360f8982dc130f18d2e9334d4dc50ed20a2a63dc8ac"
	queryAddress := "mrSoEJAs83Y46CWJikDYn7ne3Mwx3CLnkM"
	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 60000, build_tx_args, queryTxHash, queryAddress, "")
}

func test_ltc () {
	fmt.Printf("=========================\n           LTC           \n=========================\n\n")
	h := api.NewCryptocoinHandler("LTC")
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	queryTxHash := "ae1359de01c84c1750faa71ac62ed8381e97fa1156b280861ebb01fc84f538aa"
	test_common (h, nil, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_zcash () {
	fmt.Printf("=========================\n           ZCASH           \n=========================\n\n")
	h := api.NewCryptocoinHandler("ZCASH")
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	queryTxHash := "6a37436e56d4e4ac081c816b628404bc28a216afc0d0514ca0d490cc28fa5a28"
	test_common (h, nil, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_eos () {
	fmt.Printf("=========================\n           EOS           \n=========================\n\n")
	h := eos.NewEOSHandler()
	fromPrivateKey := "5JqBVZS4shWHBhcht6bn3ecWDoZXPk3TRSVpsLriQz5J3BKZtqH"
	//fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAcctName := "gzx123454321"
	//toAcctName := "eosdcrm11144"
	toAcctName := "gzxshishabi1"
	toUserKey := "dvzl22tv2jejkeq1fsgwsvd1dz2ebmlgeq"
	//queryTxHash := "a7a8fee54901cffeeb580774a867323cde87dac848d4068ca57fc4a0b4443c58"
	//queryTxHash := "e69f494b52fc59048a428856748bfd51653ea983217d79e3863f73910e64c275"
	//queryUserKey := "drefigvhv1ywn2yvjpgnu53extxb1xi4l1"

/*
	// 生成地址
	fmt.Printf("========== %s ==========\n\n", "test pubkey to address/account_name")
	address, err := h.PublicKeyToAddress(fromPubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("address from pubKeyHex is %v\n\n", address)
*/

///*
	// 构建lockin交易
	fmt.Printf("========== %s ==========\n\n", "test build unsigned transfer transaction")
	transaction, digest, err := h.BuildUnsignedLockinTransaction(fromAcctName, toUserKey, toAcctName, big.NewInt(100), "")
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

	// 组装交易
	fmt.Printf("========== %s ==========\n\n", "test make signed transaction")
	signedTransaction, err := h.MakeSignedTransaction(rsv, transaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%+v\n\n", signedTransaction)

	// 发送交易
	fmt.Printf("========== %s ==========\n\n", "test submit transaction")
	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%s\n\n", ret)
//*/

/*
	// 查交易
	fmt.Printf("========== %s ==========\n\n", "test get transaction info")
	fromAddress2, txOutputs, jsonstring, err := h.GetTransactionInfo(queryTxHash)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("from: %v\ntxOutputs: %v\njsonstring: %v\n", fromAddress2, txOutputs, jsonstring)
*/

/*
	// 查余额
	fmt.Printf("========== %s ==========\n\n", "test get balance")
	balance, err := h.GetAddressBalance(queryUserKey, "")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("balance: %v\n\n", balance)
*/
}

func test_eth () {
	fmt.Printf("=========================\n           ETH           \n=========================\n\n")
	h := api.NewCryptocoinHandler("ETH")

	//fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")

	//fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")


//fromPrivateKey, _ := crypto.HexToECDSA("0ea7b1364bc2d1d58f35324b4c0deaa129cc9bd2728e0942e7592f2836cbb530")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)
	// test 724
	fromPubKeyHex = "046377500e127f76b8403aa4e15346dbcde31c54c3914939c8216271ff159fbab5b0713615dbcf7afd262a27f4f65cbf5c5c0e41e13749fcb2956a3a008fcfc939"

	//fromAddress := "0x426B635fD6CdAf5E4e7Bf5B2A2Dd7bc6c7360FBd"
	fromAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"
	//toAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"
toAddress := "0x31341113FcDf3Eb8381D5e578CcDf06aA6720c5B"

	build_tx_args := `{"gasPrice":13000000000,"gasLimit":50000}`

	queryTxHash := "85813592d147d0e9773fcde154622ee216e16d5274f72fd8a23347300cbb667d"

//	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
queryAddress := "0xc9C0760957572F1fA90cA6Be6E43807b237C62E4"
	//queryAddress := fromAddress

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 20, build_tx_args, queryTxHash, queryAddress, "")
}

func test_etc () {
	fmt.Printf("=========================\n           ETC           \n=========================\n\n")
	h := api.NewCryptocoinHandler("ETC")

	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	queryTxHash := "0xb3c033e2f22cd8f7e6f07d3da87bd92cfe8ca6128632d6b7d596f417bc440588"

	test_common (h, fromPrivateKey, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_vechain () {
	fmt.Printf("=========================\n           VECHAIN           \n=========================\n\n")
	h := api.NewCryptocoinHandler("VECHAIN")

	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	queryTxHash := "0x1d4e7f364082341b719aa51cb86c907b37bb2335eca2b4243b6a5de39f87e87c"

	test_common (h, fromPrivateKey, fromPubKeyHex, "", "", 1, "", queryTxHash, "", "")
}

func test_erc20 () {
	fmt.Printf("=========================\n           ERC20           \n=========================\n\n")
	h := api.NewCryptocoinHandler("ERC20GUSD")

	fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
//fromPrivateKey, _ := crypto.HexToECDSA("40d6e64ce085269869b178c23a786e499ff2d6a5334fe45964211d25bea973bf")
//fromPrivateKey, _ := crypto.HexToECDSA("0ea7b1364bc2d1d58f35324b4c0deaa129cc9bd2728e0942e7592f2836cbb530")

	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	toAddress := "0x426B635fD6CdAf5E4e7Bf5B2A2Dd7bc6c7360FBd"

	build_tx_args := `"gasPrice":10000000000,"gasLimit":100000,"tokenType":"BNB"`

	queryTxHash := "0xf9e16303a1b5a59b12e18be82aaed2363621844d8b78961db57d1af7aa89419f"

//	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
	queryAddress := "0x1563E79eE9d7E4ee1246727CACabf070784F4f3b"

	query_balance_args := `"tokenType":"ERC20GUSD"`

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 20, build_tx_args, queryTxHash, queryAddress, query_balance_args)
}

// transfer at least 100000000 drops to fund a new ripple account
// 9979999990
// 79999990
func test_xrp () {
	fmt.Printf("=========================\n           XRP           \n=========================\n\n")
	h := api.NewCryptocoinHandler("XRP")
	fromKey := xrp.XRP_importKeyFromSeed("ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf", "ecdsa")
//fromKey := xrp.XRP_importKeyFromSeed("snVn7ZZnuQ5YjAGPZMD8xikMfE58z", "ecdsa")
	keyseq := uint32(0)
	fromPubKeyHex := hex.EncodeToString(fromKey.Public(&keyseq))
fmt.Printf("++++++++++++\nfromPubKeyHex is %v\n++++++++++++\n", fromPubKeyHex)
	fromAddress := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"
//fromAddress := "rJZaFWA5F4xq1rXKX8nBDL2kvZk4JazKuz"
	toAddress := "rLncvEJQhx2PY2R9X2TgBubTcpwjXm5xcs"
//toAddress := "ran6MwG2XT4N7d5d35YUhPLZK8WVoe9tnV"
//toAddress := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"
	build_tx_args := `{"fee":10}`
	fromPrivateKey := "ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf/0"
//fromPrivateKey := "snVn7ZZnuQ5YjAGPZMD8xikMfE58z/0"
	queryTxHash := "48ED82C7B3DAD0B86533B18CB5CE2BEDCE8CD841AD8930C79F428AB053FBB41C"
//queryTxHash := "50D0DA51DEB64590011D0BEDB852A811A96E5C9D3E8F162321777F31BBB30246" // lockin 100 drops
	//queryAddress := "raF1e6TSKtB34MZ9USrKphQAW5hYbARFWK"
	queryAddress := "rLncvEJQhx2PY2R9X2TgBubTcpwjXm5xcs"

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 100000000, build_tx_args, queryTxHash, queryAddress, "")
}

type Seed struct {}

func (s *Seed) Read(p []byte) (n int, err error) {
	n = 5
	return
}

func test_tron() {
	fmt.Printf("=========================\n           TRX           \n=========================\n\n")
	h := api.NewCryptocoinHandler("TRX")
	fromPrivateKey, _ := ecdsa.GenerateKey(crypto.S256(), &Seed{})
	fromPubKeyHex := trx.PublicKeyToHex(&trx.PublicKey{&fromPrivateKey.PublicKey})
	fromAddress := "417e5f4552091a69125d5dfcb7b8c2659029395bdf"
	toAddress := "41368832d8820cf4bc4f0115975e3399348c8e0249"
	//queryTxHash := "4ec77591400e30c729f10a11d6b33ef75ea565d46b6314f2fa8d170a4b8f74e1"
	queryTxHash := "3299bf491b58745020aed45ebeb139fca4c2ed86eedc57ceb336285b8416f847"
	queryAddress := fromAddress

	test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 10, "", queryTxHash, queryAddress, "")
}

func test_atom() {
	fmt.Printf("=========================\n           COSMOSATOM           \n=========================\n\n")
	h := api.NewCryptocoinHandler("ATOM")
	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)
	queryTxHash := "2A5BDD86133119F3CC3AD7BFAB986BCF8B5B7F25AA92A2A644D6B481DBBAC765"
	queryAddress := "cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt"
	test_common (h, fromPrivateKey, fromPubKeyHex, "", "", 10, "", queryTxHash, queryAddress, "")
}

func test_evt() {
	h := api.NewCryptocoinHandler("EVT1")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "EVT8JXJf7nuBEs8dZ8Pc5NpS8BJJLt6bMAmthWHE8CSqzX4VEFKtq"
	//toAddress := "EVT8PqiA7nqhaEoPXkb3zD14XK4TGugKBw3HqnkhZsvS69KrveXam"
	toAddress := "EVT5EbwKfAUyTEpQCX2U4WGf73yUmbTVGjVgrikG3Ve5ufoQyXWYc"
	//queryTxHash := "7c3c6893594568ddaae57d25d3a770e54110c8aa2b4207266cbd3db5c6cc0ac0"
	//queryTxHash := "d08a3d35cd1f31266af715535cbc2054e9bdd8799d6664126112a1b4385a514c"
	queryTxHash := "8a95dd914c8d701d992a855c447e3a828476b04e6ffcb828d0da23d0d371ea68"
	queryAddress := "EVT8JXJf7nuBEs8dZ8Pc5NpS8BJJLt6bMAmthWHE8CSqzX4VEFKtq"
        test_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, 1, "", queryTxHash, queryAddress, "")
}
