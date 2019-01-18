// 提供api的节点地址
package config

const (
	ETH_GATEWAY = "http://54.183.185.30:8018"  // Rinkeby test net
	XRP_GATEWAY = "https://s.altnet.rippletest.net:51234"  // Ripple test net
	EOS_NODEOS = "https://api.kylin.alohaeos.com"  // cryptokylin test net
	EOS_CHAIN_ID = "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191"  // cryptokylin test net
	TRON_SOLIDITY_NODE_HTTP = "https://api.shasta.trongrid.io"
)

const (
	VERSION           = 0.1
	RPCCLIENT_TIMEOUT = 30
	
	BTC_SERVER_HOST        = "47.107.50.83"  //"localhost"
	BTC_SERVER_PORT        = 8000  //18443 
	BTC_USER               = "xxmm"
	BTC_PASSWD             = "123456"
	BTC_USESSL             = false
	BTC_WALLET_PASSPHRASE  = "WalletPassphrase"
)
