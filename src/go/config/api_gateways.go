// 提供api的节点地址
package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"fmt"
	"runtime"
	"path"
	"path/filepath"
	"os"
)

type SimpleApiConfig struct {
	ApiAddress string
}

type RpcClientConfig struct {
	ElectrsAddress string
	Host string
	Port int
	User string
	Passwd string
	Usessl bool
}

type EosConfig struct {
	Nodeos string
	ChainID string
	BalanceTracker string
}

type ApiGatewayConfigs struct {
	RPCCLIENT_TIMEOUT int
	CosmosGateway *SimpleApiConfig
	TronGateway *SimpleApiConfig
	BitcoinGateway *RpcClientConfig
	OmniGateway *RpcClientConfig
	BitcoincashGateway *RpcClientConfig
	EthereumGateway *SimpleApiConfig
	EosGateway *EosConfig
	RippleGateway *SimpleApiConfig
}

var ApiGateways *ApiGatewayConfigs

var configfile string

func SetConfigFile (dir string) {
	fmt.Println("AAAAA SetConfigDir AAAAA")
	configfile = dir
}

func init () {
	fmt.Println("AAAAA init config AAAAA")
	log.Print("Loading gateway configs...")
	err := LoadApiGateways()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n================\n%+v\n================\n\n", ApiGateways)
}

func LoadApiGateways () error {
	if ApiGateways == nil {
		ApiGateways = new(ApiGatewayConfigs)
	}
	_, filename, _, _ := runtime.Caller(1)

	var configfilepath string

	configfilepath1 := configfile

	binpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configfilepath2 := path.Dir(binpath) + "/gateways.toml"

	configfilepath3 := path.Dir(filename) + "/gateways.toml"

	if exists, _ := PathExists(configfilepath1); exists {
		configfilepath = configfilepath1
	} else if exists, _ := PathExists(configfilepath2); exists {
		configfilepath = configfilepath2
	} else if exists, _ := PathExists(configfilepath3); exists {
		configfilepath = configfilepath3
	} else {
		log.Fatal(fmt.Errorf("Non of the config file path exists: %s\n%s\n%s\n"), configfilepath1, configfilepath2, configfilepath3)
	}

	_, err := toml.DecodeFile(configfilepath, ApiGateways)

	return err
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}








































// ===================== OLD CONFIGS =========================

// eth rinkeby testnet
const (
	ETH_GATEWAY = "http://54.183.185.30:8018"
)

// eos kylincrypto testnet
const (
	//eos api nodes support get actions (filter-on=*)
	EOS_NODEOS = "https://api.kylin.alohaeos.com"
	EOS_CHAIN_ID = "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191"  // cryptokylin test net
	BALANCE_SERVER = "http://127.0.0.1:7000/"
)

// ripple testnet
const (
	XRP_GATEWAY = "https://s.altnet.rippletest.net:51234"
)

// cosmos atom cosmoshub-2
var CosmosHost = "https://stargate.cosmos.network"

// tron testnet
const (
	TRON_SOLIDITY_NODE_HTTP = "https://api.shasta.trongrid.io"
)

// bitcoin testnet
const (
	ELECTRSHOST            = "http://5.189.139.168:4000"
	BTC_SERVER_HOST        = "47.107.50.83"
	BTC_SERVER_PORT        = 8000
	BTC_USER               = "xxmm"
	BTC_PASSWD             = "123456"
	BTC_USESSL             = false
)

const (
	OMNI_SERVER_HOST        = "5.189.139.168"
	OMNI_SERVER_PORT        = 9772
	OMNI_USER               = "xxmm"
	OMNI_PASSWD             = "123456"
	OMNI_USESSL             = false
)

// bitcoin cash
const (
	BCH_SERVER_HOST        = "5.189.139.168"
	BCH_SERVER_PORT        = 9552
	BCH_USER               = "xxmm"
	BCH_PASSWD             = "123456"
	BCH_USESSL             = false
)





// ===================== EVEN OLDER CONFIGS =========================

const RPCCLIENT_TIMEOUT = 30

// vechain
const (
	VECHAIN_GATEWAY = "http://127.0.0.1:50505"
)

// etc
const (
	ETC_GATEWAY = "http://127.0.0.1:50505"
)

// decred
const (
	DCR_SERVER_HOST        = "127.0.0.1"
	DCR_SERVER_PORT        = 50505
	DCR_USER               = "xxmm"
	DCR_PASSWD             = "123456"
	DCR_USESSL             = false
)

// tether
//./omnicored -conf=~/.bitcoin/bitcoin.conf -datadir=/data/usdtdata_test
//moneyqMan7uh8FqdCA2BV5yZ8qVrc9ikLP
const (
	TETHER_SERVER_HOST        = "127.0.0.1"
	TETHER_SERVER_PORT        = 50505
	TETHER_USER               = "xxmm"
	TETHER_PASSWD             = "123456"
	TETHER_USESSL             = false
)

// litecoin
const (
	LTC_SERVER_HOST        = "127.0.0.1"
	LTC_SERVER_PORT        = 50505
	LTC_USER               = "xxmm"
	LTC_PASSWD             = "123456"
	LTC_USESSL             = false
)

// bitgold
const (
	BITGOLD_SERVER_HOST        = "127.0.0.1"
	BITGOLD_SERVER_PORT        = 50505
	BITGOLD_USER               = "xxmm"
	BITGOLD_PASSWD             = "123456"
	BITGOLD_USESSL             = false
)

// dash
const (
	DASH_SERVER_HOST        = "127.0.0.1"
	DASH_SERVER_PORT        = 50505
	DASH_USER               = "xxmm"
	DASH_PASSWD             = "123456"
	DASH_USESSL             = false
)

// zcash
const (
	ZCASH_SERVER_HOST        = "127.0.0.1"
	ZCASH_SERVER_PORT        = 50505
	ZCASH_USER               = "xxmm"
	ZCASH_PASSWD             = "123456"
	ZCASH_USESSL             = false
)
