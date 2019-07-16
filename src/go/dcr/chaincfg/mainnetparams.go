// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	//"time"

	//"github.com/gaozhengxin/cryptocoins/src/go/dcr/wire"
)

var MainNetParams = Params{
	Name:        "mainnet",
	DefaultPort: "91../dcr/chaincfg/mainnetparams.go08",
	DNSSeeds: []DNSSeed{
		{"mainnet-seed.decred.mindcry.org", true},
		{"mainnet-seed.decred.netpurgatory.com", true},
		{"mainnet-seed.decred.org", true},
	},

	//GenesisHash:              &genesisHash,

	NetworkAddressPrefix: "D",
	PubKeyAddrID:         [2]byte{0x13, 0x86}, // starts with Dk
	PubKeyHashAddrID:     [2]byte{0x07, 0x3f}, // starts with Ds
	PKHEdwardsAddrID:     [2]byte{0x07, 0x1f}, // starts with De
	PKHSchnorrAddrID:     [2]byte{0x07, 0x01}, // starts with DS
	ScriptHashAddrID:     [2]byte{0x07, 0x1a}, // starts with Dc
	PrivateKeyID:         [2]byte{0x22, 0xde}, // starts with Pm
}
