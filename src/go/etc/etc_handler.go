package etc

import  (
	"context"
	"crypto/ecdsa"
	//"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"


	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gaozhengxin/cryptocoins/src/go/eth/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
	ctypes "github.com/gaozhengxin/cryptocoins/src/go/types"

)

var (
	gasPrice = big.NewInt(8000000000)
	gasLimit uint64 = 50000
	url = config.ETC_GATEWAY
	err error
	chainConfig = params.MainnetChainConfig
)

type ETCHandler struct {
}

func NewETCHandler () *ETCHandler {
	return &ETCHandler{}
}

var ETC_DEFAULT_FEE, _ = new(big.Int).SetString("10000000000",10)

func (h *ETCHandler) GetDefaultFee() *big.Int {
	return ETC_DEFAULT_FEE
}

func (h *ETCHandler) PublicKeyToAddress (pubKeyHex string) (address string, err error) {
	data := hexEncPubkey(pubKeyHex[2:])

	pub, err := decodePubkey(data)

	address = ethcrypto.PubkeyToAddress(*pub).Hex()
	return
}

// jsonstring '{"gasPrice":8000000000,"gasLimit":50000}'
func (h *ETCHandler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	var args interface{}
	json.Unmarshal([]byte(jsonstring), &args)
	json.Unmarshal([]byte(jsonstring), &args)
	if args != nil {
		userGasPrice := args.(map[string]interface{})["gasPrice"]
		userGasLimit := args.(map[string]interface{})["gasLimit"]
		if userGasPrice != nil {
			gasPrice = big.NewInt(int64(userGasPrice.(float64)))
		}
		if userGasLimit != nil {
			gasLimit = uint64(userGasLimit.(float64))
		}
	}
	transaction, hash, err := eth_newUnsignedTransaction(client, fromAddress, toAddress, amount, gasPrice, gasLimit)
	hashStr := hash.Hex()
	if hashStr[:2] == "0x" {
		hashStr = hashStr[2:]
	}
	digests = append(digests, hashStr)
	return
}

func (h *ETCHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	hashBytes, err := hex.DecodeString(hash[0])
	if err != nil {
		return
	}
	/*r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hashBytes)
	if err != nil {
		return
	}
	fmt.Printf("r: %v\ns: %v\n\n", r, s)
	rx := fmt.Sprintf("%X", r)
	sx := fmt.Sprintf("%X", s)
	rsv = append(rsv, rx + sx + "00")*/
	rsvBytes, err := ethcrypto.Sign(hashBytes, privateKey.(*ecdsa.PrivateKey))
	if err != nil {
		return
	}
	rsv = append(rsv, hex.EncodeToString(rsvBytes))
	return
}

func (h *ETCHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	return makeSignedTransaction(client, transaction.(*types.Transaction), rsv[0])
}

func (h *ETCHandler) SubmitTransaction(signedTransaction interface{}) (txhash string, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	return eth_sendTx(client, signedTransaction.(*types.Transaction))
}

func (h *ETCHandler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []ctypes.TxOutput, jsonstring string, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	hash := common.HexToHash(txhash)
	tx, isPending, err1 := client.TransactionByHash(context.Background(), hash)
	if err1 == nil && isPending == false && tx != nil {
		msg, err2 := tx.AsMessage(types.MakeSigner((*params.ChainConfig)(chainConfig), GetLastBlock()))
		err = err2
		fromAddress = msg.From().Hex()
		toAddress := msg.To().Hex()
		transferAmount := msg.Value()
		txOutput := ctypes.TxOutput{
			ToAddress: toAddress,
			Amount: transferAmount,
		}
		txOutputs = append(txOutputs, txOutput)
	} else if err1 != nil {
		err = err1
	} else if isPending {
		err = fmt.Errorf("Transaction is pending")
	} else {
		err = fmt.Errorf("Unknown error")
	}
	return
}

// args[0] coinType string
func (h *ETCHandler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	// TODO
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	account := common.HexToAddress(address)
	return client.BalanceAt(context.Background(), account, nil)
}

func GetLastBlock() *big.Int {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil
	}
	blk, _ := client.BlockByNumber(context.Background(), nil)
	return blk.Number()
}

func hexEncPubkey(h string) (ret [64]byte) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	if len(b) != len(ret) {
		panic("invalid length")
	}
	copy(ret[:], b)
	return ret
}


func decodePubkey(e [64]byte) (*ecdsa.PublicKey, error) {
	p := &ecdsa.PublicKey{Curve: ethcrypto.S256(), X: new(big.Int), Y: new(big.Int)}
	half := len(e) / 2
	p.X.SetBytes(e[:half])
	p.Y.SetBytes(e[half:])
	if !p.Curve.IsOnCurve(p.X, p.Y) {
		return nil, errors.New("invalid secp256k1 curve point")
	}
	return p, nil
}

func eth_newUnsignedTransaction (client *ethclient.Client, dcrmAddress string, toAddressHex string, amount *big.Int, gasPrice *big.Int, gasLimit uint64) (*types.Transaction, *common.Hash, error) {

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	if gasPrice == nil {
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, nil, err
		}
	}

	fromAddress := common.HexToAddress(dcrmAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, nil, err
	}

	value := amount

	toAddress := common.HexToAddress(toAddressHex)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)

	if gasLimit <= 0 {
		gasLimit, err = client.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &toAddress,
		})
		gasLimit = gasLimit * 4
		if err != nil {
			return nil, nil, err
		}
	}

	fmt.Println("gasLimit is ", gasLimit)
	fmt.Println("gasPrice is ", gasPrice)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	signer := types.NewEIP155Signer(chainID)
	txhash := signer.Hash(tx)
	return tx, &txhash, nil
}

func makeSignedTransaction(client *ethclient.Client, tx *types.Transaction, rsv string) (*types.Transaction, error) {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	message, err := hex.DecodeString(rsv)
	if err != nil {
		return nil, err
	}
	signer := types.NewEIP155Signer(chainID)
	signedtx, err := tx.WithSignature(signer, message)
	if err != nil {
		return nil, err
	}
	return signedtx, nil
}

func eth_sendTx (client *ethclient.Client, signedTx *types.Transaction) (string, error) {
	err := client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}
