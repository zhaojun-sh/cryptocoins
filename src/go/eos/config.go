package eos
import (
	eos "github.com/eoscanada/eos-go"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
)

var (
	nodeos = config.EOS_NODEOS

	opts = &eos.TxOptions{
		ChainID: hexToChecksum256(config.EOS_CHAIN_ID),
		MaxNetUsageWords: uint32(500),
		//DelaySecs: uint32(120),
		//MaxCPUUsageMS:
		Compress: eos.CompressionNone,
	}
)

const EOS_ACCURACY = 10000

const ALPHABET = "defghijklmnopqrstuvwxyz12345abcdefghijklmnopqrstuvwxyz12345abc"





