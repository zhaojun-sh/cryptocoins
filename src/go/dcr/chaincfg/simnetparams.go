// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	//"time"

)

// SimNetParams defines the network parameters for the simulation test network.
// This network is similar to the normal test network except it is intended for
// private use within a group of individuals doing simulation testing and full
// integration tests between different applications such as wallets, voting
// service providers, mining pools, block explorers, and other services that
// build on Decred.
//
// The functionality is intended to differ in that the only nodes which are
// specifically specified are used to create the network rather than following
// normal discovery rules.  This is important as otherwise it would just turn
// into another public testnet.
var SimNetParams = Params{
	Name:        "simnet",
	DefaultPort: "18555",
	DNSSeeds:    nil, // NOTE: There must NOT be any seeds.

	// Chain parameters
	//GenesisHash:              &simNetGenesisHash,
	NetworkAddressPrefix: "S",
	PubKeyAddrID:         [2]byte{0x27, 0x6f}, // starts with Sk
	PubKeyHashAddrID:     [2]byte{0x0e, 0x91}, // starts with Ss
	PKHEdwardsAddrID:     [2]byte{0x0e, 0x71}, // starts with Se
	PKHSchnorrAddrID:     [2]byte{0x0e, 0x53}, // starts with SS
	ScriptHashAddrID:     [2]byte{0x0e, 0x6c}, // starts with Sc
	PrivateKeyID:         [2]byte{0x23, 0x07}, // starts with Ps
}
