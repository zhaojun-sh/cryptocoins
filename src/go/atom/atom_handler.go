package atom

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"
	"strings"
	"github.com/btcsuite/btcd/btcec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var DefaultSendAtomFee *big.Int = big.NewInt(1)

type AtomHandler struct {}

func NewAtomHandler () *AtomHandler {
	return &AtomHandler{}
}

func (h *AtomHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error){
	pubKeyHex = strings.TrimPrefix(pubKeyHex, "0x")
	bb, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return
	}
	pk, err := btcec.ParsePubKey(bb, btcec.S256())
	if err != nil {
		return
	}
	cpk := pk.SerializeCompressed()
	var pub [33]byte
	copy(pub[:], cpk[:33])
	pubkey := secp256k1.PubKeySecp256k1(pub)
	addr := pubkey.Address()
	accAddress, err := sdk.AccAddressFromHex(addr.String())
	if err != nil {
		return
	}
	address = accAddress.String()
	return
}

func (h *AtomHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	return
}

func (h *AtomHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	return
}

func (h *AtomHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	return
}

func (h *AtomHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	return
}

func (h *AtomHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	ret, err := rpcutils.HttpGet(config.ApiGateways.CosmosGateway.ApiAddress,"txs"+"/"+txhash,nil)
	fmt.Println(string(ret))
	if err != nil {
		return
	}
	var txRes sdk.TxResponse
	err = UnmarshalJSON(ret, &txRes)
	if err != nil {
		err = errors.New("tx response error: "+string(ret)+err.Error())
		return
	}

	isSend := false
	for _, tag := range txRes.Tags {
		if tag.Key == "action" && tag.Value == "send" {
			isSend = true
			break
		}
	}
	if !isSend {
		err = errors.New("transaction does no send action")
	}

	for i, msg := range txRes.Tx.GetMsgs() {
		if txRes.Logs[i].Success {
			// confirmed = true
		}
		var atomamt sdk.Coin = sdk.NewCoin("uatom", sdk.ZeroInt())
		for _, amt := range msg.(MsgSend).Amount {
			if amt.Denom != "uatom" {
				continue
			}
			atomamt = atomamt.Add(amt)
		}
		if atomamt.Amount.Equal(sdk.ZeroInt()) {
			continue
		}
		fromAddress = msg.(MsgSend).From.String()
		toAddress := msg.(MsgSend).To.String()
		amt := atomamt.Amount.BigInt()
		txOutputs = append(txOutputs, types.TxOutput{ToAddress:toAddress,Amount:amt})
		return
	}
	return
}

func (h *AtomHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()


	ret, err := rpcutils.HttpGet(config.ApiGateways.CosmosGateway.ApiAddress,"bank/balances"+"/"+address,nil)
	if err != nil {
		return
	}
	var balances sdk.Coins
	err = json.Unmarshal(ret,&balances)
	if err != nil {
		return
	}
	for _, bal := range balances {
		if bal.Denom != "uatom" {
			continue
		}
		balance = bal.Amount.BigInt()
	}
	if balance == nil {
		balance = big.NewInt(0)
	}
	return
}

func (h *AtomHandler) GetDefaultFee () *big.Int{
	return big.NewInt(1)
}

func UnmarshalJSON (resjson []byte, res *sdk.TxResponse) (err error) {
	if res == nil {
		res = new(sdk.TxResponse)
	}
	json.Unmarshal(resjson, res)
	var resI interface{}
	json.Unmarshal(resjson, &resI)
	tx := resI.(map[string]interface{})["tx"].(map[string]interface{})
	if tx["type"] != "auth/StdTx" {
		err = errors.New("Not an auth/StdTx.")
	}
	txvalue, _ := json.Marshal(tx["value"])
	var stdTx auth.StdTx
	err = json.Unmarshal(txvalue, &stdTx)
	if &stdTx == nil {
		return
	}
	msgsI := tx["value"].(map[string]interface{})["msg"].([]interface{})
	var msgs []sdk.Msg
	for _, msgI := range msgsI {
		if msgI.(map[string]interface{})["type"] != "cosmos-sdk/MsgSend" {
			continue
		}
		msgjson, _ := json.Marshal(msgI.(map[string]interface{})["value"])
		var msg MsgSend
		json.Unmarshal(msgjson, &msg)
		msgs = append(msgs, msg)
	}
	stdTx.Msgs = msgs
	res.Tx = stdTx
	return nil
}
