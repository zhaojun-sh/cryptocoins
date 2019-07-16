package types

import "math/big"

type TxOutput struct {
	ToAddress string
	Amount *big.Int
}
