package xrp
import (
	"fmt"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"
)

func test_queryAccount () {
	key := XRP_importKeyFromSeed("shwbDeVsNa5zDYzhrbSfpxMgm1N6G", "ecdsa")
	keyseq := uint32(0)
	address := XRP_getAddress(key, &keyseq)
	fmt.Printf("address is : %v\n", address)
	account := getAccount(address)
	fmt.Printf("============ account: %v ============\n\n", account)
	seq := getSeq(address)
	fmt.Printf("============ sequence: %v ============\n\n", seq)
}

// 普通转账
func test_normalTransfer() {
	// https://developers.ripple.com/xrp-test-net-faucet.html
	seed := "snhxzrmk37AB1kkdTA8GU2Fd7iAtC"
	keyseq := uint32(0)
	XRP_Remit(seed, "ecdsa", &keyseq, "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE", big.NewInt(9979999900), int64(10))
}

func test_buildAndSendDcrmTransfer () {
	fmt.Printf("\n\n============ start ============\n\n")
	// 用dcrm公钥生成xrp地址
	dcrm_eth_address := "0x155d0ff96783CBbEc80C1Ba767A02E600f53dAEf"  // eth地址格式
	pubkeystr := "04AD1B103805DD2F41AB514CBE6C0E11EB53BEC39FF781C3A3B740867CD704AC01FD11B787F75C7511C63000CEFEA27AC10799FDAEDDF469C38CB7389B15752CAB"
	pubkey, _ := hex.DecodeString(pubkeystr)
	dcrm_xrp_address := XRP_publicKeyToAddress(pubkey)
	toaddress := "rwLc28nRV7WZiBv6vsHnpxUGAVcj8qpAtE"

/*
	// 给dcrm_xrp地址打钱激活
	XRP_FundAddress(dcrm_xrp_address)
*/

	// 从dcrm_xrp地址转出
	amount := big.NewInt(10000000)

	fmt.Printf("dcrm pubkey is %v \ndcrm xrp address is %v \ntoaddress is %v \namount is %v\n",pubkeystr, dcrm_xrp_address, toaddress, amount.String())

	tx, hash, _ := XRP_newUnsignedSimplePaymentTransaction(dcrm_xrp_address, pubkey, toaddress, amount, int64(10))

	fmt.Printf("\nunsigned transaction : %v\n\n", tx)

	// 进行dcrm签名, 暂时使用eth的api, 要使用eth地址格式
	reqstring := "{\"jsonrpc\":\"2.0\",\"method\":\"fsn_dcrmSign\",\"params\":[\"0x569882AC04C3A8831758B26B1ED343BBF7424C86CC5F8385D11C3573A811753868EA55290CCDE1B544D699CEED4DCB1BDE911DAF49B7F9D4B9480A19463D5B1F\",\"" + hash.String() + "\",\"" + dcrm_eth_address + "\",\"ETH\"],\"id\":67}"

	curlResult := DoCurlRequest(reqstring)
	rsv := getRSV(curlResult)
	sig := getSig(rsv)
	fmt.Printf("sig: %v\n\n",sig)

	// 构造交易结构, 发送交易
	XRP_makeSignedTx(tx, sig)
	res := XRP_submitTx(tx)
	fmt.Printf("%v\n",res)
}

func getRSV(curlResult string) string {
	return strings.Split(curlResult, "\"")[9]
}

func getSig(rsv string) []byte {
	l := len(rsv)-2
	rs := rsv[0:l]

	r := rs[:64]
	s := rs[64:]

	rr, _ := new(big.Int).SetString(r,16)
	ss, _ := new(big.Int).SetString(s,16)

	sign := &btcec.Signature{
		R: rr,
		S: ss,
	}
	return sign.Serialize()
}
