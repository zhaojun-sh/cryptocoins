package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"flag"
	"time"
	"github.com/gaozhengxin/cryptocoins/src/go/types"
	"github.com/gaozhengxin/cryptocoins/src/go/config"
	api "github.com/gaozhengxin/cryptocoins/src/go"
)

func main () {
	port := flag.String("port","23333","port")
	configfile := flag.String("conf","~/.gdcrm","config file")
	flag.Parse()
	path := "0.0.0.0:" + *port
	config.SetConfigFile(*configfile)
	http.HandleFunc("/gettransaction", GetTransaction)
	http.HandleFunc("/pubkeytoaddress", PubkeyToAddress)
	go http.ListenAndServe(path, nil)
	fmt.Printf("service is running on %s\n", path)
	go func () {
		for {
			log.Print("Reloading gateway configs...")
			config.LoadApiGateways()
			time.Sleep(time.Duration(60) * time.Second)
		}
	} ()
	select{}
}

type Resp struct {
	Code string `json:"code"`
	Msg string `json:"Msg,omitempty"`
	Result *GetTxResult `json:"result,omitempty"`
}

type Resp2 struct {
	Code string `json:"code"`
	Msg string `json:"Msg,omitempty"`
	Result PubkeyToAddrResult `json:"result,omitempty"`
}

var coinlist []string = []string{"BTC","ETH","ERC20BNB","ERC20HT","TRX","XRP"}

type PubkeyToAddrResult map[string]string

type GetTxResult struct {
	FromAddress string `json:"FromAddress"`
	TxOutputs []types.TxOutput `json:"TxOutputs,omitempty"`
	Err string `json:"Error,omitempty"`
}

func PubkeyToAddress (writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	pubkey, ok := request.Form["pubkey"]
	var result Resp2
	if !ok {
		result.Code = "401"
		result.Msg = "require pubkey"
	} else {
		result.Result = make(map[string]string)
		for _, cointype := range coinlist {
			h := api.NewCryptocoinHandler(cointype)
			address, err := h.PublicKeyToAddress(pubkey[0])
			if err != nil {
				result.Code = "500"
				result.Msg = err.Error()
				return
			}
			result.Result[cointype] = address
		}
		result.Code = "200"
	}
	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Fatal(err)
	}
}

func GetTransaction (writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	txhash, ok1 := request.Form["txhash"]
	cointype, ok2 := request.Form["cointype"]
	var result Resp
	if !ok1 {
		result.Code = "401"
		result.Msg = "require txhash"
	} else if !ok2 {
		result.Code = "401"
		result.Msg = "require cointype"
	} else {
		result.Code = "200"
		h := api.NewCryptocoinHandler(cointype[0])
		fromAddress, txOutputs, _, err := h.GetTransactionInfo(txhash[0])
		result.Result = &GetTxResult{
			FromAddress: fromAddress,
			TxOutputs: txOutputs,
		}
		if err != nil {
			result.Result.Err = err.Error()
		}
	}
	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Fatal(err)
	}
}
