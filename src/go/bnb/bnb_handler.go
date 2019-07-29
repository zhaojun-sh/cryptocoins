package bnb

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/client/query"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	bnbtypes  "github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var Network string = "testnet"

type BNBHandler struct {}

func NewBNBHandler () *BNBHandler {
	return &BNBHandler{}
}

func (h *BNBHandler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
	pubbytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return
	}

	pubkey, err := btcec.ParsePubKey(pubbytes, btcec.S256())
	if err != nil {
		return
	}

	cpub := pubkey.SerializeCompressed()

	pkhash := btcutil.Hash160(cpub)

	pkhashstr := hex.EncodeToString(pkhash)

	if Network == "testnet" {
		ctypes.Network = ctypes.TestNetwork
	}

	accAddress, err := ctypes.AccAddressFromHex(pkhashstr)
	if err != nil {
		return
	}

	address = accAddress.String()

	return
}

func (h *BNBHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	amt := amount.Int64()
	c := basic.NewClient("testnet-dex.binance.org:443")
	q := query.NewClient(c)
	acc, err := q.GetAccount(fromAddress)
	if err != nil {
		return
	}
	fromCoins := ctypes.Coins{{"BNB", amt}}

	fromAddr, err := ctypes.AccAddressFromBech32(fromAddress)
	if err != nil {
		return
	}

	toAddr, err := ctypes.AccAddressFromBech32(toAddress)
	if err != nil {
		return
	}

	to := []msg.Transfer{{toAddr, []ctypes.Coin{{"BNB", amt}}}}

	sendMsg := msg.CreateSendMsg(fromAddr, fromCoins, to)

	signMsg := tx.StdSignMsg{
		ChainID:"Binance-Chain-Nile",
		AccountNumber:acc.Number,
		Sequence:acc.Sequence,
		Msgs:[]msg.Msg{sendMsg},
		Memo:"this is a Dcrm lockout transaction (^_^)",
		Source:tx.Source,
	}

	transaction = BNBTx{
		SignMsg: signMsg,
		Pubkey: fromPublicKey,
	}
	digest := hex.EncodeToString(signMsg.Bytes())
	digests = append(digests, digest)
	return
}

type BNBTx struct {
	SignMsg tx.StdSignMsg
	Pubkey string
}

func (h *BNBHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	if hash[0] != "" {
		b, err := hex.DecodeString(hash[0])
		if err != nil {
			return nil, err
		}
		rs, err := privateKey.(crypto.PrivKey).Sign(b)
		if err != nil {
			return nil, err
		}
		rsv = append(rsv, hex.EncodeToString(rs)+"00")
	}
	return
}

func (h *BNBHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	if len(rsv) < 1 {
		err := fmt.Errorf("no rsv")
		return nil, err
	}
	b, err := hex.DecodeString(rsv[0])
	if err != nil {
		return
	}
	var rs []byte
	if len(b) == 65 {
		rs = b[:64]
	}
	fmt.Printf("\n\n\nlen(rs):\n%v\n\n\n", len(rs))

	signMsg := transaction.(BNBTx).SignMsg
	pubkeyhex := transaction.(BNBTx).Pubkey

	pb, err := hex.DecodeString(pubkeyhex)
	if err != nil {
		return
	}

	pub, err := btcec.ParsePubKey(pb, btcec.S256())
	if err != nil {
		return
	}

	cpub := pub.SerializeCompressed()
	var arr [33]byte
	copy(arr[:], cpub[:33])
	pubkey := secp256k1.PubKeySecp256k1(arr)

	sig := tx.StdSignature{
		AccountNumber: signMsg.AccountNumber,
		Sequence:      signMsg.Sequence,
		PubKey:        pubkey,
		Signature:     rs,
	}
	newTx := tx.NewStdTx(signMsg.Msgs, []tx.StdSignature{sig}, signMsg.Memo, signMsg.Source, signMsg.Data)
	bz, err := tx.Cdc.MarshalBinaryLengthPrefixed(&newTx)
	if err != nil {
		return
	}
	signedTransaction = []byte(hex.EncodeToString(bz))

	return
}

func (h *BNBHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	c := basic.NewClient("testnet-dex.binance.org:443")
	param := map[string]string{}
	param["sync"] = "true"
	commits, err := c.PostTx(signedTransaction.([]byte), param)
	if err != nil {
		return
	}
	txhash = commits[0].Hash
	return
}

func (h *BNBHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []types.TxOutput, jsonstring string, err error) {
	c := basic.NewClient("testnet-dex.binance.org:443")
	resp, err := c.GetTx(txhash)
	if err != nil {
		return
	}
	b, err := hex.DecodeString(resp.Data[3:len(resp.Data)-1])
	if err != nil {
		return
	}
	codec := bnbtypes.NewCodec()
	var parsedTx tx.StdTx
	err = codec.UnmarshalBinaryLengthPrefixed(b, &parsedTx)
	if err != nil {
		return
	}
	msgs := parsedTx.Msgs
	for _, m := range msgs {
		if m.Type() == "send" {
			sendmsg := m.(msg.SendMsg)
			if sendmsg.Inputs[0].Coins[0].Denom == "BNB" {
				fromAddress = sendmsg.Inputs[0].Address.String()
			}
			for _, out := range sendmsg.Outputs {
				if out.Coins[0].Denom == "BNB" {
					output := types.TxOutput{
						ToAddress: out.Address.String(),
						Amount: big.NewInt(out.Coins[0].Amount),
					}
					txOutputs = append(txOutputs, output)
				}
			}
			break
		}
	}
	return
}

func (h *BNBHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	c := basic.NewClient("testnet-dex.binance.org:443")
	q := query.NewClient(c)
	ba, err := q.GetAccount(address)
	if err != nil {
		return
	}
	for _, bal := range ba.Balances {
		var ojbk bool
		if bal.Symbol == "BNB" {
			str := strings.Replace(bal.Free.String(),".","",-1)
			balance, ojbk = new(big.Int).SetString(str, 10)
			if !ojbk {
				return nil, fmt.Errorf("parse balance error: %v", bal.Free.String())
			}
		}
	}
	return
}

func (h *BNBHandler) GetDefaultFee() *big.Int {
	return big.NewInt(50000)
}

