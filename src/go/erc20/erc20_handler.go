package erc20

import  (
	"context"
	"crypto/ecdsa"
	//"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"runtime/debug"
	"strings"

	"github.com/ethereum/go-ethereum/params"


	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"

	"github.com/gaozhengxin/cryptocoins/src/go/eth/sha3"
	"github.com/gaozhengxin/cryptocoins/src/go/erc20/token"
	ctypes "github.com/gaozhengxin/cryptocoins/src/go/types"
)

var (
	gasPrice = big.NewInt(8000000000)
	gasLimit uint64 = 50000
	url = config.ETH_GATEWAY
	err error
	chainConfig = params.RinkebyChainConfig
	//chainID = big.NewInt(40400)
)

var Tokens map[string]string = map[string]string{
	"GUSD":"0x28a79f9b0fe54a39a0ff4c10feeefa832eeceb78",
	"BNB":"0x7f30B414A814a6326d38535CA8eb7b9A62Bceae2",
	"MKR":"0x2c111ede2538400F39368f3A3F22A9ac90A496c7",
	"HT":"0x3C3d51f6BE72B265fe5a5C6326648C4E204c8B9a",
	"BNT":"0x14D5913C8396d43aB979D4B29F2102c1C65E18Db",
}

type ERC20Handler struct {
	TokenType string
}

func NewERC20Handler () *ERC20Handler {
	return &ERC20Handler{}
}

func NewERC20TokenHandler (tokenType string) *ERC20Handler {
	return &ERC20Handler{
		TokenType: tokenType,
	}
}

var ERC20_DEFAULT_FEE, _ = new(big.Int).SetString("10000000000000000",10)

func (h *ERC20Handler) GetDefaultFee() *big.Int {
	return ERC20_DEFAULT_FEE
}

func (h *ERC20Handler) PublicKeyToAddress(pubKeyHex string) (address string, err error) {
	data := hexEncPubkey(pubKeyHex[2:])

	pub, err := decodePubkey(data)

	address = ethcrypto.PubkeyToAddress(*pub).Hex()
	return
}

// jsonstring '{"gasPrice":8000000000,"gasLimit":50000,"tokenType":"BNB"}'
func (h *ERC20Handler) BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, jsonstring string) (transaction interface{}, digests []string, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	var args interface{}
	json.Unmarshal([]byte(jsonstring), &args)
	userGasPrice := args.(map[string]interface{})["gasPrice"]
	userGasLimit := args.(map[string]interface{})["gasLimit"]
	userTokenType := args.(map[string]interface{})["tokenType"]
	var tokenType string
	if userTokenType == nil {
		tokenType = h.TokenType
		if tokenType == "" {
			err = fmt.Errorf("token type not specified.")
			return
		}
		if Tokens[tokenType] == "" {
			err = fmt.Errorf("token not supported")
			return
		}
	}
	if userGasPrice != nil {
		gasPrice = big.NewInt(int64(userGasPrice.(float64)))
	}
	if userGasLimit != nil {
		gasLimit = uint64(userGasLimit.(float64))
	}
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	transaction, hash, err := erc20_newUnsignedTransaction(client, fromAddress, toAddress, amount, gasPrice, gasLimit, tokenType)
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

func (h *ERC20Handler) SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error) {
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

func (h *ERC20Handler) MakeSignedTransaction(rsv []string, transaction interface{}) (signedTransaction interface{}, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	return makeSignedTransaction(client, transaction.(*types.Transaction), rsv[0])
}

func (h *ERC20Handler) SubmitTransaction(signedTransaction interface{}) (ret string, err error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	return erc20_sendTx(client, signedTransaction.(*types.Transaction))
}

func (h *ERC20Handler) GetTransactionInfo(txhash string) (fromAddress string, txOutputs []ctypes.TxOutput, jsonstring string, err error) {
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
		data := msg.Data()

		toAddress, transferAmount, decodeErr := DecodeTransferData(data)
		txOutput := ctypes.TxOutput{
			ToAddress: toAddress,
			Amount: transferAmount,
		}
		txOutputs = append(txOutputs, txOutput)
		if decodeErr != nil {
			err = decodeErr
			return
		}
	} else if err1 != nil {
		err = err1
	} else if isPending {
		err = fmt.Errorf("Transaction is pending")
	} else {
		err = fmt.Errorf("Unknown error")
	}
	return
}

// jsonstring:'{"tokenType":"BNB"}'
func (h *ERC20Handler) GetAddressBalance(address string, jsonstring string) (balance *big.Int, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	var args interface{}
	json.Unmarshal([]byte(jsonstring), &args)
	tokenType := args.(map[string]interface{})["tokenType"]
	if tokenType != nil {
		err = fmt.Errorf("token type not specified")
	}

	tokenAddr := Tokens[tokenType.(string)]
	if tokenAddr == "" {
		err = fmt.Errorf("Token not supported")
		return
	}

	myABIJson := `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"}]`
	myABI, err := abi.JSON(strings.NewReader(myABIJson))
	if err != nil {
		return
	}

	data, err := myABI.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		return
	}
	dataHex := "0x" + hex.EncodeToString(data)
	fmt.Printf("data is %v\n\n", dataHex)

	reqJson := `{"jsonrpc": "2.0","method": "eth_call","params": [{"to": "` + tokenAddr + `","data": "` + dataHex + `"},"latest"],"id": 1}`
	fmt.Printf("reqJson: %v\n\n", reqJson)

	ret := rpcutils.DoPostRequest2(url, reqJson)
	fmt.Printf("ret: %v\n\n", ret)

	var retStruct map[string]interface{}
	json.Unmarshal([]byte(ret), &retStruct)
	if retStruct["result"] == nil {
		if retStruct["error"] != nil {
			err = fmt.Errorf(retStruct["error"].(map[string]interface{})["message"].(string))
			return
		}
		err = fmt.Errorf(ret)
		return
	}
	balanceStr := retStruct["result"].(string)[2:]
	balanceHex, _ := new(big.Int).SetString(balanceStr, 16)
	balance, _ = new(big.Int).SetString(fmt.Sprintf("%d",balanceHex), 10)

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress(tokenAddr)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	balance1, _ := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	fmt.Printf("balance1: %v\n\n", balance1)

	return
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

func DecodeTransferData(data []byte) (toAddress string, transferAmount *big.Int, err error) {
	eventData := data[:4]
	if string(eventData) == string([]byte{0xa9, 0x05, 0x9c, 0xbb}) {
		addressData := data[4:36]
		amountData := data[36:]
		num, _ := new(big.Int).SetString(hex.EncodeToString(addressData), 16)
		toAddress = "0x" + fmt.Sprintf("%x", num)
		amountHex, _ := new(big.Int).SetString(hex.EncodeToString(amountData), 16)
		transferAmount, _ = new(big.Int).SetString(fmt.Sprintf("%d", amountHex), 10)
	} else {
		err = fmt.Errorf("Invalid transfer data")
		return
	}
	return
}

func erc20_newUnsignedTransaction (client *ethclient.Client, dcrmAddress string, toAddressHex string, amount *big.Int, gasPrice *big.Int, gasLimit uint64, tokenType string) (*types.Transaction, *common.Hash, error) {

	chainID, err := client.NetworkID(context.Background())

	if err != nil {
		return nil, nil, err
	}

	tokenAddressHex, ok := Tokens[tokenType]
	if ok {
	} else {
		err = errors.New("token not supported")
		return nil, nil, err
	}

	if gasPrice == nil {
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, nil, err
		}
	}

	fromAddress := common.HexToAddress(dcrmAddress)
/*
	nonce or pending nonce
*/
	//nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, nil, err
	}

	value := big.NewInt(0)

	toAddress := common.HexToAddress(toAddressHex)
	tokenAddress := common.HexToAddress(tokenAddressHex)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	if gasLimit <= 0 {
		gasLimit, err = client.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: data,
		})
		gasLimit = gasLimit * 4
		if err != nil {
			return nil, nil, err
		}
	}

	fmt.Println("gasLimit is ", gasLimit)
	fmt.Println("gasPrice is ", gasPrice)
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

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

func erc20_sendTx (client *ethclient.Client, signedTx *types.Transaction) (string, error) {
	err := client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}
