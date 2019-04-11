package cryptocoins

import (
	"regexp"
)

var RegExpmap map[string]string = map[string]string {
	"BTC":"^(1|3|m|n)[a-zA-Z\\d]{24,33}$",
	"ETH":"^(0x)?[0-9a-fA-F]{40}$",
	"XRP":"^r[1-9a-km-zA-HJ-NP-Z]{33}$",
}

type AddressValidator struct {
	Exp string
}

func NewAddressValidator (cointype string) *AddressValidator {
	return &AddressValidator{
		Exp: RegExpmap[cointype],
	}
}

func (v *AddressValidator) IsValidAddress (address string) bool {
	match, _ := regexp.MatchString(v.Exp, address)
	return match
}
