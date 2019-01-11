# 方法列表
| 方法签名 | 描述 |
|:----------:|-------------|
|PublicKeyToAddress (pubKeyHex string) (address, address21 string, err error)|通过公钥生成地址  |
|BuildUnsignedTransaction(fromAddress, fromPublicKey, toAddress string, amount *big.Int, args ...interface{}) (transaction interface{}, digests []string, err error)|构建未签名交易  |
|SignTransaction(hash []string, privateKey interface{}) (rsv []string, err error)|用私钥对交易哈希签名  |
|MakeSignedTransaction (rsv []string, transaction interface{}) (signedTransaction interface{}, err error)|构建带签名的交易  |
|SubmitTransaction(signedTransaction interface{}) (ret string, err error)|发送交易到网络  |
|GetTransactionInfo(txhash string) (fromAddress, toAddress string, transferAmount *big.Int, err error)|根据交易哈希查询交易信息  |
|GetAddressBalance(address string) (balance *big.Int, err error)|查询地址余额  |

# 运行demo程序
$ cd ./src/go/trx  
$ go run *.go  
