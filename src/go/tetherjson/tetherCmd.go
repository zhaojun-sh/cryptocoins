package tetherjson

import (
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
)

type OmniGetTransactionCmd struct {
	Txid string
}

func NewOmniGetTransactionCmd(txHash string) *OmniGetTransactionCmd {
	return &OmniGetTransactionCmd{
		Txid: txHash,
	}
}

func init () {
	err := btcjson.RegisterCmd("omni_gettransaction", (*OmniGetTransactionCmd)(nil), 0x1)
	fmt.Println(err)
}
