package eos
import (
	eos "github.com/eoscanada/eos-go"

	"github.com/gaozhengxin/cryptocoins/src/go/config"
)

var (
	nodeos = config.EOS_NODEOS

	opts = &eos.TxOptions{
		ChainID: hexToChecksum256(config.EOS_CHAIN_ID),
		MaxNetUsageWords: uint32(999),
		//DelaySecs: uint32(120),
		MaxCPUUsageMS: uint8(200),
		Compress: eos.CompressionNone,
	}
)

const CREATOR_ACCOUNT = "gzx123454321"

const CREATOR_PRIVKEY = "5JqBVZS4shWHBhcht6bn3ecWDoZXPk3TRSVpsLriQz5J3BKZtqH"

const EOS_ACCURACY = 10000

const ALPHABET = "defghijklmnopqrstuvwxyz12345abcdefghijklmnopqrstuvwxyz12345abc"

const BALANCE_SERVER = "http://127.0.0.1:7000/"

var InitialRam = uint32(256)

var InitialCPU = int64(1000)

var InitialStakeNet = int64(1000)



