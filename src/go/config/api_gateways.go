// 提供api的节点地址
package config

// eth rinkeby testnet
const (
	ETH_GATEWAY = "http://54.183.185.30:8018"
	//ETH_GATEWAY = "https://rinkeby.infura.io"
	//ETH_GATEWAY = "http://127.0.0.1:8111"
)

// etc
const (
	ETC_GATEWAY = "http://127.0.0.1:50505"
)

// vechain
const (
	VECHAIN_GATEWAY = "http://127.0.0.1:50505"
)

// eos kylincrypto testnet
const (
	//eos api nodes support get actions (filter-on=*)
	EOS_NODEOS = "https://api.kylin.alohaeos.com"
	//EOS_NODEOS = "https://api-kylin.eoslaomao.com"
	EOS_CHAIN_ID = "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191"  // cryptokylin test net
)

// ripple testnet
const (
	XRP_GATEWAY = "https://s.altnet.rippletest.net:51234"
)

// tron testnet
const (
	TRON_SOLIDITY_NODE_HTTP = "https://api.shasta.trongrid.io"
)

// bitcoin testnet
const (
	BTC_SERVER_HOST        = "47.107.50.83"
	BTC_SERVER_PORT        = 8000
	BTC_USER               = "xxmm"
	BTC_PASSWD             = "123456"
	BTC_USESSL             = false
	BTC_WALLET_PASSPHRASE  = "WalletPassphrase"
)

// decred
const (
	RPCCLIENT_TIMEOUT = 30
	DCR_SERVER_HOST        = "127.0.0.1"
	DCR_SERVER_PORT        = 50505
	DCR_USER               = "xxmm"
	DCR_PASSWD             = "123456"
	DCR_USESSL             = false
	DCR_WALLET_PASSPHRASE  = "WalletPassphrase"
)

// tether
const (
	TETHER_SERVER_HOST        = "127.0.0.1"
	TETHER_SERVER_PORT        = 50505
	TETHER_USER               = "xxmm"
	TETHER_PASSWD             = "123456"
	TETHER_USESSL             = false
	TETHER_WALLET_PASSPHRASE  = "WalletPassphrase"
)

// litecoin
const (
	LTC_SERVER_HOST        = "127.0.0.1"
	LTC_SERVER_PORT        = 50505
	LTC_USER               = "xxmm"
	LTC_PASSWD             = "123456"
	LTC_USESSL             = false
	LTC_WALLET_PASSPHRASE  = "WalletPassphrase"
)

// bitcoin cash
const (
	BCH_SERVER_HOST        = "127.0.0.1"
	BCH_SERVER_PORT        = 50607
	BCH_USER               = "xxmm"
	BCH_PASSWD             = "123456"
	BCH_USESSL             = false
	BCH_WALLET_PASSPHRASE  = "WalletPassphrase"
)

// dash
const (
	DASH_SERVER_HOST        = "127.0.0.1"
	DASH_SERVER_PORT        = 50505
	DASH_USER               = "xxmm"
	DASH_PASSWD             = "123456"
	DASH_USESSL             = false
	DASH_WALLET_PASSPHRASE  = "WalletPassphrase"
)

// zcash
const (
	ZCASH_SERVER_HOST        = "127.0.0.1"
	ZCASH_SERVER_PORT        = 50505
	ZCASH_USER               = "xxmm"
	ZCASH_PASSWD             = "123456"
	ZCASH_USESSL             = false
	ZCASH_WALLET_PASSPHRASE  = "WalletPassphrase"
)
