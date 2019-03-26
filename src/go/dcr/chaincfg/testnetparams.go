// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	//"time"

)

// TestNet3Params defines the network parameters for the test currency network.
// This network is sometimes simply called "testnet".
// This is the third public iteration of testnet.
var TestNet3Params = Params{
	Name:        "testnet3",
	DefaultPort: "19108",
	DNSSeeds: []DNSSeed{
		{"testnet-seed.decred.mindcry.org", true},
		{"testnet-seed.decred.netpurgatory.com", true},
		{"testnet-seed.decred.org", true},
	},

	// Chain parameters
	//GenesisHash:              &testNet3GenesisHash,
	NetworkAddressPrefix: "T",
	PubKeyAddrID:         [2]byte{0x28, 0xf7}, // starts with Tk
	PubKeyHashAddrID:     [2]byte{0x0f, 0x21}, // starts with Ts
	PKHEdwardsAddrID:     [2]byte{0x0f, 0x01}, // starts with Te
	PKHSchnorrAddrID:     [2]byte{0x0e, 0xe3}, // starts with TS
	ScriptHashAddrID:     [2]byte{0x0e, 0xfc}, // starts with Tc
	PrivateKeyID:         [2]byte{0x23, 0x0e}, // starts with Pt
}
