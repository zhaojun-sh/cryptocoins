package main

import (
	"fmt"
	"github.com/ellsol/evt/evtapi/client"
	//"github.com/ellsol/evt/evtapi/v1/evt"
	"github.com/ellsol/evt/evtapi/v1/chain"
	//"github.com/ellsol/evt/evtapi/v1/history"
	"github.com/ellsol/evt/evtconfig"
	"github.com/sirupsen/logrus"
)

var evtcfg *evtconfig.Instance

func init () {
	evtcfg = evtconfig.New("https://testnet1.everitoken.io")
}

func main () {
	clt := client.New(evtcfg, logrus.New())
	//apievt := evt.New(evtcfg, clt)
	apichain := chain.New(evtcfg, clt)
	//apihistory := history.New(evtcfg, clt)
/*
	fmt.Printf("\n\n========== evt/get_fungible ============\n\n")
	res1, err := apievt.GetFungible("1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n\n", res1)
*/
/*
	fmt.Printf("\n\n========== chain/get_block_header_state ============\n\n")
	res2, err := apichain.GetHeadBlockHeaderState()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n\n", res2)
*/
	fmt.Printf("\n\n========== chain/get_block ============\n\n")
	res3, err := apichain.GetBlock("00000337011f960e705815a53fd7525d7bd7caab8292aa5962a3f63770d1d0ba")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n\n", res3)
/*
	fmt.Printf("\n\n========== chain/get_transaction_ids_for_block ============\n\n")
	res4, err := apichain.GetTransactionIdsForBlock("00000337011f960e705815a53fd7525d7bd7caab8292aa5962a3f63770d1d0ba")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n\n", res4)

	fmt.Printf("\n\n========== history/get_transaction ============\n\n")
	res5, err := apihistory.GetTransaction("b30c14ea6708c5fc5e83cee9876f45519eb0895131f2fe3330c7ce1210f81093")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n\n", res5)
*/
}
