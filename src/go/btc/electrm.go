package btc

import (
	"fmt"
	"encoding/json"
	rpcutils "github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
	"runtime/debug"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	"github.com/fusion/go-fusion/log"
	"sort"
)

func listUnspent_electrs(addr string) (list []btcjson.ListUnspentResult, err error) {
	defer func () {
		if e := recover(); e != nil {
			err = fmt.Errorf("Runtime error: %v\n%v", e, string(debug.Stack()))
			return
		}
	} ()
	path := `address/` + addr + `/utxo`
	ret, err := rpcutils.HttpGet(config.ApiGateways.BitcoinGateway.ElectrsAddress, path, nil)
	if err != nil {
		return
	}
	var utxos []electrsUtxo
	err = json.Unmarshal(ret, &utxos)
	if err != nil {
		return
	}
	fmt.Printf("\n\n%v\n\n", string(ret))
	fmt.Printf("\n\n%+v\n\n", utxos)
	for _, utxo := range utxos {
		path = `tx/` + utxo.Txid
		txret, txerr := rpcutils.HttpGet(config.ApiGateways.BitcoinGateway.ElectrsAddress, path, nil)
		if txerr != nil {
			log.Debug("======== get utxo script ========", "error", txerr)
			continue
		}
		var tx electrsTx
		txerr = json.Unmarshal(txret, &tx)
		if txerr != nil {
log.Debug("======== get utxo script ========", "error", txerr)
			continue
		}
		utxo.Script = tx.Vout[int(utxo.Vout)].Scriptpubkey
		res := btcjson.ListUnspentResult{
			TxID: utxo.Txid,
			Vout: uint32(utxo.Vout),
			ScriptPubKey: utxo.Script,
			Address: addr,
			Amount: utxo.Value/1e8,
			Spendable: true,
		}
		if utxo.Status.Confirmed {
			res.Confirmations = 6
		} else {
			res.Confirmations = 0
		}
		list = append(list, res)
	}
	sort.Sort(sortableLURSlice(list))
log.Debug("======== get utxo ========", "utxo list", list)
	return
}

type electrsTx struct {
	Txid string
	Vout []electrsTxOut
}

type electrsTxOut struct {
	Scriptpubkey string
}

type electrsUtxo struct {
	Txid string `json:"txid"`
	Vout uint32
	Script string
	Status utxoStatus
	Value float64
}

type utxoStatus struct {
	Confirmed bool
	Block_height float64
	Block_hash string
	Block_time float64
}
