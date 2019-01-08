package main

import (
	"fmt"
	"math/big"
)

var pubKeyHex = "04E657EE43DDADBA5AF1F1AD4E8098D996C2F3C397E807C9255B0850EA2D151D050F4A0D6451231DD35F0FF653F166C35BCAA0E817520B4DB87DE7E060A72D578E"

func main() {
	test_erc20()
}

func test_erc20 () {
	h := NewTransactionHandler("ERC20")
	fmt.Printf("h type: %T\n\n", h)

	address, err := h.PublicKeyToAddress(pubKeyHex)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("address is %v\n\n", address)

	address1 := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"
	address2 := "0x7b5Ec4975b5fB2AA06CB60D0187563481bcb6140"

	transaction, digest, err := h.BuildUnsignedTransaction(address1, address2, big.NewInt(1), big.NewInt(10), uint64(5000000), "GUSD")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n%v\n\n", transaction, digest)

	var rsv []string
	rsv = append(rsv, "26C729F4B7C1D0407BB0B8D1052771B20ED3DC96739EDC2694684EF3FCA935735B38E837D6BAE347BFE2B29772211C3FE66164B72B1FA23F612A4B9148D013C500")
	signedTransaction, err := h.MakeSignedTransaction(rsv, transaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n", signedTransaction)

	ret, err := h.SubmitTransaction(signedTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%s\n\n", ret)
}
