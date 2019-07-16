package omni

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type RpcResult struct {
	Result OmniTx `json:"result"`
	Error string `json:"error"`
}

type OmniTx struct {
	Confirmations int64 `json:"confirmations"`
	Valid bool `json:"valid"`
	From string `json:"sendingaddress"`
	To string `json:"referenceaddress"`
	AmountString string `json:"amount"`
	Amount *big.Int
	Type string `json:"type"`
	PropertyName string `json:"propertyname"`
	Error error
}

func DecodeOmniTx(ret string) (*OmniTx) {
	var res  = new(RpcResult)
	json.Unmarshal([]byte(ret), res)
	omnitx := res.Result
	if res.Error != "" {
		omnitx.Error = fmt.Errorf(res.Error)
	}
	s := strings.Replace(omnitx.AmountString,".","",-1)
	omnitx.Amount, _ = new(big.Int).SetString(s,10)
	omnitx.PropertyName = "OMNI"+omnitx.PropertyName
	return &omnitx
}
