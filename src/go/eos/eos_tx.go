package eos
import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/rubblelabs/ripple/crypto"

	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
)

type EOSTransactionHandler struct {
}

func (h *EOSTransactionHandler) PublicKeyToAddress(pubKeyHex string) (acctName string, msg string, err error) {
	acctName = GenRandomAccountName(pubKeyHex)
	msg = "this account name is not expected to be seen on chain, and the public key is not delegated by any existed account. it is possible to create an account with this name and delegate the public key."
	pubKey, err := HexToPubKey(pubKeyHex)
	if err != nil {
		return
	}
	accounts, err := GetAccountNameByPubKey(pubKey.String())
	if len(accounts) != 0 {
		msg = "this public key is delegated by accounts: "
		for _, acct := range accounts {
			msg = msg + "/" + acct
		}
	}
	return
}

// args[0] memo string
func (h *EOSTransactionHandler) BuildUnsignedTransaction(fromAcctName, fromPublicKey, toAcctName string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error) {
	if fromAcctName == "" {
		accounts, err1 := GetAccountNameByPubKey(fromPublicKey)
		if err1 != nil {
			err = err1
			return
		}
		if len(accounts) > 1 {
			err = fmt.Errorf("public key is delegated by multiple accounts, need to specify account name")
		}
		fromAcctName = accounts[0]
	}
	memo := args[0].(string)
	digest, transaction, err := EOS_newUnsignedTransaction(fromAcctName, toAcctName, amount, memo)
	if err != nil {
		return
	}

	digests = append(digests, digest)

	return
}

func (h *EOSTransactionHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	signature, err := SignDigestWithPrivKey(hash[0], privateKey.(string))
	if err != nil {
		return
	}
	vrs := signature.Content
	rsvBytes := append(vrs[1:], vrs[0])
	rsv = append(rsv, hex.EncodeToString(rsvBytes))
	return
}

func (h *EOSTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	signature, err := RSVToSignature(rsv[0])
	if err != nil {
		return
	}
	signedTransaction = MakeSignedTransaction(transaction.(*eos.SignedTransaction), signature)
	return
}

func (h *EOSTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	ret = SubmitTransaction(signedTransaction.(*eos.SignedTransaction))
	return
}

func (h *EOSTransactionHandler) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	api := "v1/history/get_transaction"
	data := `{"id":"` + txhash + `","block_num_hint":"0"}`
	ret := rpcutils.DoCurlRequest(nodeos, api, data)
	var retStruct map[string]interface{}
	json.Unmarshal([]byte(ret), &retStruct)
	if retStruct["trx"] == nil {
		if reterr := retStruct["error"]; reterr != nil {
			name := reterr.(map[string]interface{})["name"]
			details := reterr.(map[string]interface{})["details"].([]interface{})
			var message string
			if details != nil {
				message = details[0].(map[string]interface{})["message"].(string)
			}
			err = fmt.Errorf("%v, message: %v", name, message)
			return
		}
		err = fmt.Errorf("  %v", ret)
		return
	}
	tfData := retStruct["trx"].(map[string]interface{})["actions"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})
	fromAddress = tfData["from"].(string)
	toAddress = tfData["receiver"].(string)
	transferAmount = big.NewInt(int64(tfData["transfer"].(float64)))
	return
}

func (h *EOSTransactionHandler) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error) {
	api := "v1/chain/get_account"
	data := `{"account_name":"` + address + `"}`
	ret := rpcutils.DoCurlRequest(nodeos, api, data)
	var retStruct map[string]interface{}
	json.Unmarshal([]byte(ret), &retStruct)
	if retStruct["core_liquid_balance"] == nil {
		err = fmt.Errorf("%v", ret)
		return
	}
	balStr := retStruct["core_liquid_balance"].(string)
	fmt.Printf("%s\n\n", balStr)
	fmt.Printf("%s\n\n", strings.Fields(balStr)[0])
	balFloat, _ := strconv.ParseFloat(strings.Fields(balStr)[0], 64)
	balInt := int64(balFloat * EOS_ACCURACY)
	balance = big.NewInt(balInt)
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
		return
	}
}

func checkAPIErr(res string) error {
	var v interface{}
	err := json.Unmarshal([]byte(res), &v)
	if err != nil {
		return err
	}
	m := v.(map[string]interface{})
	if m["error"] != nil {
		if errm := m["error"].(map[string]interface{}); len(errm) > 0 {
			err = fmt.Errorf(errm["name"].(string), errm)
			return err
		}
	}
	return nil
}

func isCanonical(compactSig []byte) bool {
	// From EOS's codebase, our way of doing Canonical sigs.
	// https://steemit.com/steem/@dantheman/steem-and-bitshares-cryptographic-security-update
	//
	// !(c.data[1] & 0x80)
	// && !(c.data[1] == 0 && !(c.data[2] & 0x80))
	// && !(c.data[33] & 0x80)
	// && !(c.data[33] == 0 && !(c.data[34] & 0x80));

	d := compactSig
	t1 := (d[1] & 0x80) == 0
	t2 := !(d[1] == 0 && ((d[2] & 0x80) == 0))
	t3 := (d[33] & 0x80) == 0
	t4 := !(d[33] == 0 && ((d[34] & 0x80) == 0))

	return t1 && t2 && t3 && t4
}

func GetHeadBlockID(nodeos string) (chainID string, err error) {
	api := "v1/chain/get_info"
	res := rpcutils.DoCurlRequest(nodeos, api, "")
	if err = checkAPIErr(res); err != nil {
		return "", err
	}
	var v interface{}
	json.Unmarshal([]byte(res), &v)
	m := v.(map[string]interface{})
	return fmt.Sprintf("%v",m["head_block_id"]), nil
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

func SignDigestWithPrivKey(hash, wif string) (ecc.Signature, error) {
	digest := hexToChecksum256(hash)
	privKey, err := ecc.NewPrivateKey(wif)
	checkErr(err)
	return privKey.Sign(digest)
}

func GetAccountNameByPubKey(pubKey string) ([]string, error) {
	api := "v1/history/get_key_accounts"
	data := "{\"public_key\":\"" + pubKey + "\"}"
	res := rpcutils.DoCurlRequest(nodeos, api, data)
	if err := checkAPIErr(res); err != nil {
		return nil, err
	}
	var v interface{}
	json.Unmarshal([]byte(res), &v)
	m := v.(map[string]interface{})
	accts := m["account_names"].([]interface{})
	var accounts []string
	for _, acct := range accts {
		accounts = append(accounts, acct.(string))
	}
	return accounts, nil
}

func RSVToSignature (rsvStr string) (ecc.Signature, error) {
	rsv, _ := hex.DecodeString(rsvStr)
	rsv[64] += byte(31)
	v := rsv[64]
	rs := rsv[:64]
	vrs := append([]byte{v}, rs...)
	data := append([]byte{0}, vrs...)
	return ecc.NewSignatureFromData(data)
}

func hexToChecksum256(data string) eos.Checksum256 {
	bytes, err := hex.DecodeString(data)
	checkErr(err)
	return eos.Checksum256(bytes)
}

// 根据公钥生成随机的账户名
func GenRandomAccountName(pubKeyHex string) string {
	b, _ := hex.DecodeString(pubKeyHex)
	fmt.Printf("!!! %v \n!!! %v\n\n", pubKeyHex, b)

	r, _ := rand.Int(rand.Reader, big.NewInt(256))
	b = append(b, byte(r.Uint64()))

	b = btcutil.Hash160(b)

	b = append([]byte{0}, b...)
	return crypto.Base58Encode(b, ALPHABET)[:12]
}

func EOS_newUnsignedTransaction(fromAcctName, toAcctName string, amount *big.Int, memo string) (string, *eos.SignedTransaction, error) {

	from := eos.AccountName(fromAcctName)
	to := eos.AccountName(toAcctName)
	s := strconv.FormatFloat(float64(amount.Int64())/10000, 'f', 4, 64) + " EOS"
	quantity, _ := eos.NewAsset(s)

	transfer := &eos.Action{
		Account: eos.AN("eosio.token"),
		Name:    eos.ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{
				Actor: from,
				Permission: eos.PN("active"),
			},
		},
		ActionData: eos.NewActionData(token.Transfer{
			From:     from,
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}

        var actions []*eos.Action
        actions = append(actions, transfer)

	// 获取 head block id
	hbid, err := GetHeadBlockID(nodeos)
	if err != nil {
		return "", nil, err
	}
	opts.HeadBlockID = hexToChecksum256(hbid)
        tx := eos.NewTransaction(actions, opts)

	stx := eos.NewSignedTransaction(tx)

	txdata, cfd, err := stx.PackedTransactionAndCFD()
	checkErr(err)
	digest := eos.SigDigest(opts.ChainID, txdata, cfd)
	digestStr := hex.EncodeToString(digest)
	return digestStr, stx, nil
}

func MakeSignedTransaction(stx *eos.SignedTransaction, signature ecc.Signature) *eos.SignedTransaction {
	stx.Signatures = append(stx.Signatures, signature)
	return stx
}

func SubmitTransaction (stx *eos.SignedTransaction) string {

	txjson := stx.String()

	b := "{\"signatures\":[\"" + stx.Signatures[0].String() + "\"], \"compression\":\"none\", \"transaction\":" + txjson + "}"

	res := rpcutils.DoCurlRequest(nodeos, "v1/chain/push_transaction", b)
	return res
}

// 创建eos账户
// 需要一个creator账户, creator要有余额用于购买内存
func CreateNewAccount(creatorName, creatorActivePrivKey, accountName, accountPubKey string, buyram uint32) (string, error) {

	opubKey, err := ecc.NewPublicKey(accountPubKey)
	if err != nil {
		return "", err
	}
	apubKey, err := ecc.NewPublicKey(accountPubKey)
	if err != nil {
		return "", err
	}

	// 创建账户action
	action1 := &eos.Action{
		Account: eos.AccountName("eosio"),
		Name:    eos.ActionName("newaccount"),
		Authorization: []eos.PermissionLevel{
			{eos.AccountName(creatorName), eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(system.NewAccount{
			Creator: eos.AccountName(creatorName),
			Name:    eos.AccountName(accountName),
			Owner: eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: opubKey,
						Weight:    1,
					},
				},
			},
			Active: eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: apubKey,
						Weight:    1,
					},
				},
			},
		}),
	}

	// 买内存action
	action2 := system.NewBuyRAMBytes(eos.AccountName(creatorName), eos.AccountName(accountName), buyram)

	// 获取 head block id
	hbid, err := GetHeadBlockID(nodeos)
	checkErr(err)
	opts.HeadBlockID = hexToChecksum256(hbid)

	// 创建账户和买内存一定要同时执行
	actions := []*eos.Action{action1, action2}

	tx := eos.NewTransaction(actions, opts)

	stx := eos.NewSignedTransaction(tx)

	txdata, cfd, err := stx.PackedTransactionAndCFD()
        checkErr(err)
        digest := eos.SigDigest(opts.ChainID, txdata, cfd)
        digestStr := hex.EncodeToString(digest)

        signature, err := SignDigestWithPrivKey(digestStr, creatorActivePrivKey)
	if err != nil {
		return "", err
	}

        stx.Signatures = append(stx.Signatures, signature)

        txjson := stx.String()

        b := "{\"signatures\":[\"" + stx.Signatures[0].String() + "\"], \"compression\":\"none\", \"transaction\":" + txjson + "}"

        res := rpcutils.DoCurlRequest(nodeos, "v1/chain/push_transaction", b)
	if err = checkAPIErr(res); err != nil {
		return "", err
	}

	return res, nil
}

// 预购cpu和net带宽, 用于帐号执行各种action
func DelegateBW (fromAcctName, fromActivePrivKey, receiverName string, stakeCPU, stakeNet int64, transfer bool) (string, error){
	from := eos.AccountName(fromAcctName,)
	receiver := eos.AccountName(receiverName)
	action := system.NewDelegateBW(from, receiver, eos.NewEOSAsset(stakeCPU), eos.NewEOSAsset(stakeNet), transfer)

	// 获取 head block id
	hbid, err := GetHeadBlockID(nodeos)
	checkErr(err)
	opts.HeadBlockID = hexToChecksum256(hbid)

	actions := []*eos.Action{action}
	tx := eos.NewTransaction(actions, opts)
	stx := eos.NewSignedTransaction(tx)
	txdata, cfd, err := stx.PackedTransactionAndCFD()
	checkErr(err)
	digest := eos.SigDigest(opts.ChainID, txdata, cfd)
	digestStr := hex.EncodeToString(digest)
	signature, err := SignDigestWithPrivKey(digestStr, fromActivePrivKey)
	if err != nil {
		return "", err
	}
	stx.Signatures = append(stx.Signatures, signature)

	txjson := stx.String()
	b := "'{\"signatures\":[\"" + stx.Signatures[0].String() + "\"], \"compression\":\"none\", \"transaction\":" + txjson + "}'"
	res := rpcutils.DoCurlRequest(nodeos, "v1/chain/push_transaction", b)
	if err = checkAPIErr(res); err != nil {
		return "", err
	}
	return res, nil
}
