package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"strings"
	"log"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/gaozhengxin/cryptocoins/src/go/xrp"
	"github.com/gaozhengxin/cryptocoins/src/go/trx"
	api "github.com/gaozhengxin/cryptocoins/src/go"
)

func mustnotnil (name string, obj interface{}) {
	if obj == nil {
		log.Fatal(name + " is nil.")
	}
}

func main () {
	to := flag.String("to","","to address")
	ct := flag.String("cointype","","coin type")
	amount := flag.String("amount","","amount")
	flag.Parse()
	mustnotnil("to address", to)
	mustnotnil("cointype", ct)
	mustnotnil("amount", amount)
	toAddress := *to
	cointype := *ct
	amt, ok := new(big.Int).SetString(*amount, 10)
	if !ok {
		log.Fatal("invalid amount.")
	}
	sender := GetSender(cointype)
	mustnotnil("sender", sender)
	sender(toAddress, amt)
}

func GetSender (coinType string) func (toAddress string, amt *big.Int) {
	coinTypeC := strings.ToUpper(coinType)
	switch coinTypeC {
	case "BTC":
		return func(toAddress string, amt *big.Int) {send_btc(toAddress, amt)}
	case "ETH":
		return func(toAddress string, amt *big.Int) {send_eth(toAddress, amt)}
	case "TRX":
		return func(toAddress string, amt *big.Int) {send_tron(toAddress, amt)}
	case "XRP":
		return func(toAddress string, amt *big.Int) {send_xrp(toAddress, amt)}
	default:
		if isErc20(coinTypeC) {
			return func(toAddress string, amt *big.Int) {send_erc20(toAddress, amt, coinTypeC)}
		}
		if isEVT(coinTypeC) {
			return func(toAddress string, amt *big.Int) {send_evt(toAddress, amt, coinTypeC)}
		}
		return nil
	}
}

func isErc20(tokentype string) bool {
	return strings.HasPrefix(tokentype,"ERC20")
}

func isEVT(tokentype string) bool {
	return strings.HasPrefix(tokentype,"EVT")
}

func send_common (h api.CryptocoinHandler, fromPrivateKey interface{}, fromPubKeyHex, fromAddress, toAddress string, amt *big.Int, build_tx_args string, queryTxHash, queryAddress string, query_balance_args string) {
	fmt.Printf("========== %s ==========\n\n", "build unsigned transfer transaction")
	transaction, digest, err := h.BuildUnsignedTransaction(fromAddress, fromPubKeyHex, toAddress, amt, build_tx_args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("transaction: %+v\n\ndigest: %v\n\n", transaction, digest)

	fmt.Printf("========== %s ==========\n\n", "sign with private key")
	rsv, err := h.SignTransaction(digest, fromPrivateKey)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("rsv is %+v\n\n", rsv)

	fmt.Printf("========== %s ==========\n\n", "make signed transaction")
	signedTransaction, err := h.MakeSignedTransaction(rsv, transaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%+v\n\n", signedTransaction)

	fmt.Printf("========== %s ==========\n\n", "submit transaction")
	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
	}
	fmt.Printf("%s\n\n", ret)

}

func send_btc (toAddress string, amt *big.Int) {
	fmt.Printf("=========================\n           BTC           \n=========================\n\n")
	h := api.NewCryptocoinHandler("BTC")
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "mtjq9RmBBDVne7YB4AFHYCZFn3P2AXv9D5"
	//fromAddress := "mv6iMFM84xVLh7s6tR1ryD6PkmMeJEkkuh"
//toAddress := "moneyqMan7uh8FqdCA2BV5yZ8qVrc9ikLP"
	build_tx_args := `{"feeRate":0.0001}`
	queryTxHash := "c89e489d0368a498537892e45f8825ee683f7e144bcbd6b8891c9eac0ba01807"
	queryAddress := "mrSoEJAs83Y46CWJikDYn7ne3Mwx3CLnkM"
	send_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, amt, build_tx_args, queryTxHash, queryAddress, "")
}

func send_eth (toAddress string, amt *big.Int) {
	fmt.Printf("=========================\n           ETH           \n=========================\n\n")
	h := api.NewCryptocoinHandler("ETH")

	//fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
	fromPrivateKey, _ := crypto.HexToECDSA("d55b502bd4867b2c1b505af9b7cefeeb910b6cfbb570e2e47680bc89ee123eab")
//fromPrivateKey, _ := crypto.HexToECDSA("0ea7b1364bc2d1d58f35324b4c0deaa129cc9bd2728e0942e7592f2836cbb530")
	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x426B635fD6CdAf5E4e7Bf5B2A2Dd7bc6c7360FBd"
	//toAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	build_tx_args := `{"gasPrice":8000000000,"gasLimit":50000}`

	queryTxHash := "85813592d147d0e9773fcde154622ee216e16d5274f72fd8a23347300cbb667d"

//	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
queryAddress := "0xc9C0760957572F1fA90cA6Be6E43807b237C62E4"
	//queryAddress := fromAddress

	send_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, amt, build_tx_args, queryTxHash, queryAddress, "")
}

func send_erc20 (toAddress string, amt *big.Int, tokentype string) {
	fmt.Printf("=========================\n           ERC20           \n=========================\n\n")
	h := api.NewCryptocoinHandler(tokentype)

	fromPrivateKey, _ := crypto.HexToECDSA("a751c37b0a6e4b7605512fefb28cd4bd141bc3c06863557624800140eddf13be")
//fromPrivateKey, _ := crypto.HexToECDSA("40d6e64ce085269869b178c23a786e499ff2d6a5334fe45964211d25bea973bf")
//fromPrivateKey, _ := crypto.HexToECDSA("0ea7b1364bc2d1d58f35324b4c0deaa129cc9bd2728e0942e7592f2836cbb530")

	pub := crypto.FromECDSAPub(&fromPrivateKey.PublicKey)
	fromPubKeyHex := hex.EncodeToString(pub)

	fromAddress := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"


	build_tx_args := `"gasPrice":20000000000,"gasLimit":100000,"tokenType":"BNB"`

	queryTxHash := "0xf9e16303a1b5a59b12e18be82aaed2363621844d8b78961db57d1af7aa89419f"

//	queryAddress := "0xEc430068f392e5FBcE92016758C5111375d16f7D"
	queryAddress := "0x1563E79eE9d7E4ee1246727CACabf070784F4f3b"

	query_balance_args := `"tokenType":"ERC20GUSD"`

	send_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, amt, build_tx_args, queryTxHash, queryAddress, query_balance_args)
}

// transfer at least 100000000 drops to fund a new ripple account
// 9979999990
// 79999990
func send_xrp (toAddress string, amt *big.Int) {
	fmt.Printf("=========================\n           XRP           \n=========================\n\n")
	h := api.NewCryptocoinHandler("XRP")
	fromKey := xrp.XRP_importKeyFromSeed("ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf", "ecdsa")
//fromKey := xrp.XRP_importKeyFromSeed("snVn7ZZnuQ5YjAGPZMD8xikMfE58z", "ecdsa")
	keyseq := uint32(0)
	fromPubKeyHex := hex.EncodeToString(fromKey.Public(&keyseq))
fmt.Printf("++++++++++++\nfromPubKeyHex is %v\n++++++++++++\n", fromPubKeyHex)
	fromAddress := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"
//fromAddress := "rJZaFWA5F4xq1rXKX8nBDL2kvZk4JazKuz"
//toAddress := "ran6MwG2XT4N7d5d35YUhPLZK8WVoe9tnV"
//toAddress := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"
	build_tx_args := `{"fee":10}`
	fromPrivateKey := "ssfL5tmpTTqCw5sHjnRHQ4yyUCQKf/0"
//fromPrivateKey := "snVn7ZZnuQ5YjAGPZMD8xikMfE58z/0"
	queryTxHash := "48ED82C7B3DAD0B86533B18CB5CE2BEDCE8CD841AD8930C79F428AB053FBB41C"
//queryTxHash := "50D0DA51DEB64590011D0BEDB852A811A96E5C9D3E8F162321777F31BBB30246" // lockin 100 drops
	//queryAddress := "raF1e6TSKtB34MZ9USrKphQAW5hYbARFWK"
	queryAddress := "rLncvEJQhx2PY2R9X2TgBubTcpwjXm5xcs"

	send_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, amt, build_tx_args, queryTxHash, queryAddress, "")
}

type Seed struct {}

func (s *Seed) Read(p []byte) (n int, err error) {
	n = 5
	return
}

func send_tron(toAddress string, amt *big.Int) {
	fmt.Printf("=========================\n           TRX           \n=========================\n\n")
	h := api.NewCryptocoinHandler("TRX")
	fromPrivateKey, _ := ecdsa.GenerateKey(crypto.S256(), &Seed{})
	fromPubKeyHex := trx.PublicKeyToHex(&trx.PublicKey{&fromPrivateKey.PublicKey})
	fromAddress := "417e5f4552091a69125d5dfcb7b8c2659029395bdf"
	//queryTxHash := "4ec77591400e30c729f10a11d6b33ef75ea565d46b6314f2fa8d170a4b8f74e1"
	queryTxHash := "3299bf491b58745020aed45ebeb139fca4c2ed86eedc57ceb336285b8416f847"
	queryAddress := fromAddress

	send_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, amt, "", queryTxHash, queryAddress, "")
}

func send_evt(toAddress string, amt *big.Int, tokentype string) {
	fmt.Printf("=========================\n           EVT           \n=========================\n\n")
	h := api.NewCryptocoinHandler(tokentype)
	fromPrivateKey := "93N2nFzgr1cPRU8ppswy8HrgBMaoba8aH5sGZn9NdgG9weRFrA1"
	fromPubKeyHex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	fromAddress := "EVT8JXJf7nuBEs8dZ8Pc5NpS8BJJLt6bMAmthWHE8CSqzX4VEFKtq"
	send_common (h, fromPrivateKey, fromPubKeyHex, fromAddress, toAddress, amt, "", "", "", "")
}
