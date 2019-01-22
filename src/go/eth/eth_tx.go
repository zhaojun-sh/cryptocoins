package eth

import  (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"


	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gaozhengxin/cryptocoins/src/go/config"

)

var (
	url = config.ETH_GATEWAY
	err error
	chainConfig = params.RinkebyChainConfig
)

type ETHTransactionHandler struct {
}

func (h *ETHTransactionHandler) PublicKeyToAddress (pubKeyHex string) (address string, msg string, err error) {
	data := hexEncPubkey(pubKeyHex[2:])

	pub, err := decodePubkey(data)

	address = ethcrypto.PubkeyToAddress(*pub).Hex()
	return
}

//args[0]: gasPrice	*big.Int
//args[1]: gasLimit	uint64
//args[2]: tokenType	string
func (h *ETHTransactionHandler) BuildUnsignedTransaction (fromAddress, fromPublicKey, toAddress string, amount *big.Int, args []interface{}) (transaction interface{}, digests []string, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	transaction, hash, err := eth_newUnsignedTransaction(client, fromAddress, toAddress, amount, args[0].(*big.Int), args[1].(uint64))
	hashStr := hash.Hex()
	if hashStr[:2] == "0x" {
		hashStr = hashStr[2:]
	}
	digests = append(digests, hashStr)
	return
}

/*
func SignTransaction(hash string, address string) (rsv string, err error) {
	return
}
*/

func (h *ETHTransactionHandler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
	hashBytes, err := hex.DecodeString(hash[0])
	if err != nil {
		return
	}
	r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hashBytes)
	if err != nil {
		return
	}
	fmt.Printf("r: %v\ns: %v\n\n", r, s)
	rx := fmt.Sprintf("%X", r)
	sx := fmt.Sprintf("%X", s)
	rsv = append(rsv, rx + sx + "00")
	return
}

func (h *ETHTransactionHandler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	return makeSignedTransaction(client, transaction.(*types.Transaction), rsv[0])
}

func (h *ETHTransactionHandler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	return eth_sendTx(client, signedTransaction.(*types.Transaction))
}

func (h *ETHTransactionHandler) GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, _ []interface{}, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	hash := common.HexToHash(txhash)
	tx, isPending, err1 := client.TransactionByHash(context.Background(), hash)
	if err1 == nil && isPending == false && tx != nil {
		msg, err2 := tx.AsMessage(types.MakeSigner(chainConfig, GetLastBlock()))
		err = err2
		fromAddress = msg.From().Hex()
		toAddress = msg.To().Hex()
		transferAmount = msg.Value()
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
func (h *ETHTransactionHandler) GetAddressBalance(address string, args []interface{}) (balance *big.Int, err error) {
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
	return "success/" + signedTx.Hash().Hex(), nil
}
