// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	//"math"
	//"time"

)

// RegNetParams defines the network parameters for the regression test network.
// This should not be confused with the public test network or the simulation
// test network.  The purpose of this network is primarily for unit tests and
// RPC server tests.  On the other hand, the simulation test network is intended
// for full integration tests between different applications such as wallets,
// voting service providers, mining pools, block explorers, and other services
// that build on Decred.
//
// Since this network is only intended for unit testing, its values are subject
// to change even if it would cause a hard fork.
var RegNetParams = Params{
	Name:        "regnet",
	DefaultPort: "18655",
	DNSSeeds:    nil, // NOTE: There must NOT be any seeds.

	// Chain parameters
	//GenesisHash:              &regNetGenesisHash,

	NetworkAddressPrefix: "R",
	PubKeyAddrID:         [2]byte{0x25, 0xe5}, // starts with Rk
	PubKeyHashAddrID:     [2]byte{0x0e, 0x00}, // starts with Rs
	PKHEdwardsAddrID:     [2]byte{0x0d, 0xe0}, // starts with Re
	PKHSchnorrAddrID:     [2]byte{0x0d, 0xc2}, // starts with RS
	ScriptHashAddrID:     [2]byte{0x0d, 0xdb}, // starts with Rc
	PrivateKeyID:         [2]byte{0x22, 0xfe}, // starts with Pr
}
