// 提供api的节点地址
package config

// eth rinkeby testnet
const (
	ETH_GATEWAY = "http://54.183.185.30:8018"
	//ETH_GATEWAY = "https://rinkeby.infura.io"
	//ETH_GATEWAY = "http://127.0.0.1:8111"
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
	XRP_GATEWAY = "https://s.altnet.rippletest.net:51234"  // Ripple test net
)

// tron testnet
const (
	TRON_SOLIDITY_NODE_HTTP = "https://api.shasta.trongrid.io"
)

// bitcoin testnet
const (
	RPCCLIENT_TIMEOUT = 30
	BTC_SERVER_HOST        = "47.107.50.83"
	BTC_SERVER_PORT        = 8000 
	BTC_USER               = "xxmm"
	BTC_PASSWD             = "123456"
	BTC_USESSL             = false
	BTC_WALLET_PASSPHRASE  = "WalletPassphrase"
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

// bitcash
const (
	BCH_SERVER_HOST        = "127.0.0.1"
	BCH_SERVER_PORT        = 50607 
	BCH_USER               = "xxmm"
	BCH_PASSWD             = "123456"
	BCH_USESSL             = false
	BCH_WALLET_PASSPHRASE  = "WalletPassphrase"
)
