package main

import (
	"fmt"
	"github.com/gaozhengxin/cryptocoins/src/go/eos"
)

/*func main() {
	accountName := "shiwodadiao0"
	ok1, e1 := eos.CreateNewAccount(creatorName, creatorActivePrivKey, accountName, ownerkey, activekey, eos.InitialRam)
	ok2, e2 := eos.DelegateBW(creatorName, creatorActivePrivKey, accountName, eos.InitialCPU, eos.InitialCPU, true)
}*/

func main() {
	accountName := "gzx123454321"
	creatorName := "gzx123454321"
	creatorActivePrivKey := "5JqBVZS4shWHBhcht6bn3ecWDoZXPk3TRSVpsLriQz5J3BKZtqH"
	ok, e := eos.BuyRAM(creatorName, creatorActivePrivKey, accountName, eos.InitialRam*4)
	fmt.Printf("%v, %v\n",ok,e)
}
