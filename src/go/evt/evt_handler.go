package evt

import (
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"github.com/gaozhengxin/cryptocoins/src/go/eos"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/ellsol/evt/ecc"
	"github.com/ellsol/evt/evtapi/client"
	"github.com/ellsol/evt/evtapi/v1/evt"
	"github.com/ellsol/evt/evtapi/v1/chain"
	"github.com/ellsol/evt/evtapi/v1/history"
	"github.com/ellsol/evt/evtconfig"
	"github.com/ellsol/evt/evttypes"
	//"github.com/ellsol/evt/transaction/fungible"
	//ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/sirupsen/logrus"
)

type EvtHandler struct {
	TokenId uint
}

var r *rand.Rand

// 只支持fungible token
func NewEvtHandler (tokenId string) *EvtHandler {
	tid, err := strconv.Atoi(strings.TrimPrefix(tokenId,"EVT"))
	if err != nil {
		return nil
	}
	return &EvtHandler{
		TokenId: uint(tid),
	}
}

// EVT地址就是EVT格式的pubkey
func (h *EvtHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error){
	pk, err := HexToPubKey(pubKeyHex)
	if err != nil {
		return
	}
	address = pk.String()
	return
}

func (h *EvtHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {

	key := strconv.Itoa(int(h.TokenId))
	number := "0.00010 S#1"

	// 1. abi_json_to_bin https://www.everitoken.io/developers/apis,_sdks_and_tools/abi_reference
	args := chain.Args{
		Transfer:chain.ActionType{
			Name:"transfer",
			Threshold:1,
			Authorizers:[]chain.Authorizers{chain.Authorizers{Ref:"[A] "+fromAddress,Weight:1}},
		},
		From:fromAddress,
		To:toAddress,
		Number:number,
		Memo:"this is a dcrm lockout (^_^)",
	}
	actarg := chain.ActionArguments{
		Action:"transferft",
		Args:args,
	}
	bb, _ := json.Marshal(actarg)
	fmt.Printf("\n%+v\n",string(bb))

	evtcfg := evtconfig.New(config.ApiGateways.EVTGateway.ApiAddress)
	clt := client.New(evtcfg, logrus.New())
	apichain := chain.New(evtcfg, clt)

	res, apierr := apichain.AbiJsonToBin(&actarg)
	if apierr != nil {
		err = apierr.Error()
		return
	}

	// 2. evttypes.Trxjson
	action := evttypes.Action{
		Name:"transferft",
		Domain:".fungible",
		Key:key,
	}
	trx := &evttypes.TRXJson{
		MaxCharge: 10000,
		Actions: []evttypes.SimpleAction{evttypes.SimpleAction{Action:action,Data:res.Binargs}},
		Payer: fromAddress,
		TransactionExtensions: make([]interface{},0),
	}

	// 3. chain/trx_json_to_digest expiration??? ref_block_num??? ref_block_prefix??? ...
	layout := "2006-01-02T15:04:05"

	res2, apierr := apichain.GetInfo()
	if apierr != nil {
		err = apierr.Error()
		return
	}
	fmt.Printf("\n\ncnm\ngetinfo result\n%+v\nmsln\n\n",res2)

	headtime, _ := time.Parse(layout,res2.HeadBlockTime)
	exptime := headtime.Add(time.Duration(60)*time.Minute)

	trx.Expiration = exptime.Format(layout)

	trx.RefBlockNum = res2.LastIrreversibleBlockNum
	//trx.RefBlockNum = res2.HeadBlockNum

	res3, apierr := apichain.GetBlock (strconv.Itoa(trx.RefBlockNum))
	if apierr != nil {
		err = apierr.Error()
		return
	}
	trx.RefBlockPrefix = res3.RefBlockPrefix

	b, _ := json.Marshal(trx)
	fmt.Printf("\ncnm\ntrx is \n%+v\n\n%v\nnmsl\n\n",trx,string(b))

	// 4. TRXJsonToDigest
	res4, apierr := apichain.TRXJsonToDigest(trx)
	if apierr != nil {
		err = apierr.Error()
		return
	}
	fmt.Printf("\ncnm\nnmsl\nres is \n%+v\n\n",res3)

	transaction = trx
	digests = append(digests,res4.Digest)
	return
}

func (h *EvtHandler) SignTransaction (hash []string, wif interface{}) (rsv []string, err error) {
/*
	// 直接用ecdsa签名, 可能报is not canonical
	pkwif, err :=  btcutil.DecodeWIF(wif.(string))
	if err != nil {
		return
	}
	privateKey := pkwif.PrivKey
	hashBytes, err := hex.DecodeString(hash[0])
	if err != nil {
		return
	}
	rsvBytes, err := ethcrypto.Sign(hashBytes, privateKey.ToECDSA())
	if err != nil {
		return
	}
	rsv = append(rsv, hex.EncodeToString(rsvBytes))
*/


///*
	var rsvBytes []byte
	for i:=0; i<25; i++ {
		fmt.Printf("\ntry %v\n", i)
		src := rand.NewSource(time.Now().UnixNano())
		r = rand.New(src)
		pkwif, err1 :=  btcutil.DecodeWIF(wif.(string))
		if err1 != nil {
			err = err1
			return
		}
		privateKey := pkwif.PrivKey
		hashBytes, err2 := hex.DecodeString(hash[0])
		if err2 != nil {
			err = err2
			return
		}
		//rsvBytes1, err3 := ethcrypto.Sign(hashBytes, privateKey.ToECDSA())
		sig, err3 := RandNonceSign(privateKey, hashBytes)
		if err3 != nil {
			err = err3
			return
		}
		rsvBytes1 := append(sig.R.Bytes(), sig.S.Bytes()...)
		pk := privateKey.PubKey().SerializeUncompressed()
		v := byte(0)
		fmt.Printf("msg length is %v\n",len(hashBytes))
		for j := 0; j < (btcec.S256().H+1)*2; j ++ {
			rsvBytes2 := append(rsvBytes1, byte(j))
			pkr, e := secp256k1.RecoverPubkey(hashBytes,rsvBytes2)
			if e == nil && Equal(pk, pkr) {
				fmt.Printf("v = %v, ojbk\n",j)
				v = byte(j)
				break
			}
			fmt.Printf("%v\n",e)
			fmt.Printf("pk is %v\npkr is %v\n", pk, pkr)
			fmt.Printf("v = %v, not ok\n",j)
		}
		rsvBytes1 = append(rsvBytes1, v)
		fmt.Printf("\nrsvBytes1 is %v\n\n", rsvBytes1)
		compactSig, err4 := eos.RSVToSignature(hex.EncodeToString(rsvBytes1))
		if err4 != nil {
			err = err4
			return
		}
		fmt.Printf("%v\n", compactSig)
		if eos.IsCanonical(compactSig.Content) {
			rsvBytes = rsvBytes1
			fmt.Println("is canonical!")
			break
		}
	}
	if len(rsvBytes) == 0 {
		return nil, fmt.Errorf("!!!!\nfail to produce a canonical signature\n!!!!")
	}
	rsv = append(rsv, hex.EncodeToString(rsvBytes))
//*/


	// 用eoscanada/eos-go的方法签名, 不会报is not canonical
/*
	privateKey, err := ecc.NewPrivateKey(wif.(string))
	if err != nil {
		return
	}
	digest := eos.HexToChecksum256(hash[0])
	sig, err := privateKey.Sign(digest)
	if err != nil {
		return
	}
	vrs := sig.Content
	v := vrs[0] - byte(31)
	rsvBytes := append(vrs[1:], v)
	rsv = append(rsv, hex.EncodeToString(rsvBytes))
*/

	fmt.Printf("!!!!!!!! Sign Transaction !!!!!!!!!    rsv is %v\n\n",rsv)
	return
}

func Equal(x []byte, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i, xi := range x {
		if xi != y[i] {
			return false
		}
	}
	return true
}

func RandNonceSign(privateKey *btcec.PrivateKey, hash []byte) (*btcec.Signature, error){
	if r == nil {
		src := rand.NewSource(time.Now().UnixNano())
		r = rand.New(src)
	}
	privkey := privateKey.ToECDSA()
	N := btcec.S256().N
	halfOrder := new(big.Int).Rsh(N, 1)
	// 搞一个随机数
	// k := nonceRFC6979(privkey.D, hash)
	k := new(big.Int).Add(new(big.Int).Rand(r, new(big.Int).Sub(N,big.NewInt(1000))),big.NewInt(1000))

	inv := new(big.Int).ModInverse(k, N)
	r, _ := privkey.Curve.ScalarBaseMult(k.Bytes())
	r.Mod(r, N)

	if r.Sign() == 0 {
		return nil, fmt.Errorf("calculated R is zero")
	}

	e := hashToInt(hash, privkey.Curve)
	s := new(big.Int).Mul(privkey.D, r)
	s.Add(s, e)
	s.Mul(s, inv)
	s.Mod(s, N)

	if s.Cmp(halfOrder) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, fmt.Errorf("calculated S is zero")
	}
	return &btcec.Signature{R: r, S: s}, nil
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

func (h *EvtHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	sig, err := eos.RSVToSignature(rsv[0])
	if err != nil {
		return
	}
	// evttypes.SignedTRXJson
	signedTransaction = &evttypes.SignedTRXJson{
		Signatures: []string{sig.String()},
		Compression: "none",
		Transaction: transaction.(*evttypes.TRXJson),
	}
	return
}

func (h *EvtHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	// chain/push_transaction
	evtcfg := evtconfig.New(config.ApiGateways.EVTGateway.ApiAddress)
	clt := client.New(evtcfg, logrus.New())
	apichain := chain.New(evtcfg, clt)
	b, _ := json.Marshal(signedTransaction)
	fmt.Println(string(b))
	res, apierr := apichain.PushTransaction(signedTransaction.(*evttypes.SignedTRXJson))
	if apierr != nil {
		err = apierr.Error()
		return
	}
	txhash = res.TransactionId
	return
}

func (h *EvtHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()

	evtcfg := evtconfig.New(config.ApiGateways.EVTGateway.ApiAddress)
	clt := client.New(evtcfg, logrus.New())
	apihistory := history.New(evtcfg, clt)
	res, apierr := apihistory.GetTransaction(txhash)
	if apierr != nil {
		err = apierr.Error()
		return
	}
	fromAddress = res.InnerTransaction.Payer
	actions := res.InnerTransaction.Actions
	var transfer *chain.Action
	for _, act := range actions {
		if act.Name == "transferft" || act.Name == "issuefungible" {
			transfer = &act
			txout, err := parseAction(h.TokenId, transfer)
			if err != nil {
				continue
			}
			txOutputs = append(txOutputs, *txout)
		}
	}
	return
}

func parseAction (tarid uint, transfer *chain.Action) (*types.TxOutput, error) {
	if transfer.Name == "transferft" {
		tmp := strings.Split(transfer.Data.Number,"#")[1]
		symid, _ := strconv.Atoi(tmp)
		if uint(symid) != tarid {
			return nil, fmt.Errorf("sym id is %v, want %v", symid, tarid)
		}
		amtstr := strings.Replace(strings.Split(transfer.Data.Number," ")[0],".","",-1)
		fmt.Printf("amtstr is %s\n", amtstr)
		amt, ok := new(big.Int).SetString(amtstr, 10)
		if !ok {
			err := fmt.Errorf("transfer amount error: %s", transfer.Data.Number)
			return nil, err
		}
		txout := &types.TxOutput{ToAddress:transfer.Data.To,Amount:amt}
		fmt.Printf("txout is %+v\n", txout)
		return txout, nil
	}
	if transfer.Name == "issuefungible" {
		amtstr := strings.Replace(strings.Split(transfer.Data.Number," ")[0],".","",-1)
                amt, ok := new(big.Int).SetString(amtstr, 10)
                if !ok {
			err := fmt.Errorf("transfer amount error: %s", transfer.Data.Number)
                        return nil, err
                }
		txout := &types.TxOutput{ToAddress:transfer.Data.Address,Amount:amt}
		return txout, nil
	}
	return nil, fmt.Errorf("evt parse action: unknown error.")
}

func (h *EvtHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()

	evtcfg := evtconfig.New(config.ApiGateways.EVTGateway.ApiAddress)
	clt := client.New(evtcfg, logrus.New())
	apievt := evt.New(evtcfg, clt)
	res, apierr := apievt.GetFungibleBalance(h.TokenId, address)
	if apierr != nil {
		err = apierr.Error()
		return
	}
	amtstr := strings.Replace(strings.Split((*res)[0]," ")[0],".","",-1)
	balance, ok := new(big.Int).SetString(amtstr, 10)
	if !ok {
		err = fmt.Errorf("transfer amount error: %s", (*res)[0])
		return
	}
	return
}

func (h *EvtHandler) GetDefaultFee () *big.Int{
	return big.NewInt(1)
}

func PubKeyToHex(pk string) (pubKeyHex string, _ error) {
	pubKey, err := ecc.NewPublicKey(pk)
	if err != nil {
		return "", err
	}
	pubKeyHex = "0x" + hex.EncodeToString(pubKey.Content)
	return
}

func HexToPubKey(pubKeyHex string) (ecc.PublicKey, error) {
	fmt.Printf("hex is %v\nlen(hex) is %v\n\n", pubKeyHex, len(pubKeyHex))
	if pubKeyHex[:2] == "0x" || pubKeyHex[:2] == "0X" {
		pubKeyHex = pubKeyHex[2:]
	}
	// TODO 判断长度
	if len(pubKeyHex) == 130 {
		uBytes, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			return ecc.PublicKey{}, err
		}
		pubkey, err := btcec.ParsePubKey(uBytes, btcec.S256())
		if err != nil {
			return ecc.PublicKey{}, err
		}
		pubkeyBytes := pubkey.SerializeCompressed()
		pubkeyBytes = append([]byte{0}, pubkeyBytes...)  // byte{0} 表示 curve K1, byte{1} 表示 curve R1
		return ecc.NewPublicKeyFromData(pubkeyBytes)
	}

	if len(pubKeyHex) == 66 {
		pubkeyBytes, _ := hex.DecodeString(pubKeyHex)
		pubkeyBytes = append([]byte{0}, pubkeyBytes...)
		return ecc.NewPublicKeyFromData(pubkeyBytes)
	}

	return ecc.PublicKey{}, fmt.Errorf("unexpected public key length  %v", len(pubKeyHex))
}
