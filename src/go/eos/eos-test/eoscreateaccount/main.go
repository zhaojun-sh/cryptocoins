package main

import (
	//"fmt"
	"github.com/gaozhengxin/cryptocoins/src/go/eos"
)

func main() {
	//accountName := "shiwodadiao0"
	accountName := "woshishabi00"
	creatorName := "gzx111111111"
	creatorActivePrivKey := "5HuM1Co9E9vUmfUnXhyduAx3seaCQd1ginmnmtV37fHd6BFgFPY"
	ok1, e1 := eos.CreateNewAccount(creatorName, creatorActivePrivKey, accountName, ownerkey, activekey, eos.InitialRam)
	ok2, e2 := eos.DelegateBW(creatorName, creatorActivePrivKey, accountName, eos.InitialCPU, eos.InitialStakeNet, true)
}

/*
func main() {
//	accountName := "gzx111111111"
//	creatorName := "gzx123454321"
//	creatorActivePrivKey := "5JqBVZS4shWHBhcht6bn3ecWDoZXPk3TRSVpsLriQz5J3BKZtqH"

	accountName := "gzx123454321"
	creatorName := "gzx111111111"
	creatorActivePrivKey := "5HuM1Co9E9vUmfUnXhyduAx3seaCQd1ginmnmtV37fHd6BFgFPY"

	// 可以给自己买ram
	ok1, e1 := eos.BuyRAM(creatorName, creatorActivePrivKey, accountName, eos.InitialRam*10)
	fmt.Printf("%v, %v\n",ok1,e1)
	// 不能给自己买cpu和stake net
	ok2, e2 := eos.DelegateBW(creatorName, creatorActivePrivKey, accountName, eos.InitialCPU*4, eos.InitialStakeNet, true)
	fmt.Printf("%v, %v\n",ok2,e2)
}
*/
