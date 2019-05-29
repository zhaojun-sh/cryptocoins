package evt

import (
	"fmt"
	"math/big"
	"encoding/hex"
	"runtime/debug"
	"strconv"
	"strings"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ellsol/evt/ecc"
	"github.com/ellsol/evt/evtapi/client"
	"github.com/ellsol/evt/evtapi/v1/evt"
	"github.com/ellsol/evt/evtapi/v1/chain"
	"github.com/ellsol/evt/evtapi/v1/history"
	"github.com/ellsol/evt/evtconfig"
	//"github.com/ellsol/evt/evttypes"
	"github.com/sirupsen/logrus"
)

type EvtHandler struct {
	TokenId uint
}

// 只支持fungible token
func NewEvtHandler (tokenId string) *EvtHandler {
	tid, err := strconv.Atoi(strings.TrimPrefix(tokenId,"EVT-"))
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

	// 1. fungible.TransferFungibleParams.Arguments() -> evttypes.Trxjson
	// 2. abi_json_to_bin https://www.everitoken.io/developers/apis,_sdks_and_tools/abi_reference
	// 3. chain/trx_json_to_digest expiration??? ref_block_num??? ref_block_prefix??? ...

	return
}

func (h *EvtHandler) SignTransaction (hash []string, privateKey interface{}) (rsv []string, err error) {
	return
}

func (h *EvtHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	// evttypes.SignedTRXJson
	return
}

func (h *EvtHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	// chain/push_transaction
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
