package main

import (
	"fmt"
	"github.com/gaozhengxin/cryptocoins/src/go/eos"
)


func main () {
	pubhex := "04c1a8dd2d6acd8891bddfc02bc4970a0569756ed19a2ed75515fa458e8cf979fdef6ebc5946e90a30c3ee2c1fadf4580edb1a57ad356efd7ce3f5c13c9bb4c78f"
	pubkey, err := eos.HexToPubKey(pubhex)
	fmt.Printf("%v, %v",pubkey.String(), err)
}
