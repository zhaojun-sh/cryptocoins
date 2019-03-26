// Copyright (c) 2013, 2014 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcrutil

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ripemd160"

	"github.com/gaozhengxin/cryptocoins/src/go/dcr/base58"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/chaincfg"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/chaincfg/chainec"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrec"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrec/edwards"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrec/secp256k1"
	"github.com/gaozhengxin/cryptocoins/src/go/dcr/dcrec/secp256k1/schnorr"
)

var (
	ErrChecksumMismatch = errors.New("checksum mismatch")

	ErrUnknownAddressType = errors.New("unknown address type")

	ErrAddressCollision = errors.New("address collision")

	ErrMissingDefaultNet = errors.New("default net not defined")
)

func encodeAddress(hash160 []byte, netID [2]byte) string {
	return base58.CheckEncode(hash160[:ripemd160.Size], netID)
}

func encodePKAddress(serializedPK []byte, netID [2]byte, algo dcrec.SignatureType) string {
	pubKeyBytes := []byte{0x00}

	switch algo {
	case dcrec.STEcdsaSecp256k1:
		pubKeyBytes[0] = byte(dcrec.STEcdsaSecp256k1)
	case dcrec.STEd25519:
		pubKeyBytes[0] = byte(dcrec.STEd25519)
	case dcrec.STSchnorrSecp256k1:
		pubKeyBytes[0] = byte(dcrec.STSchnorrSecp256k1)
	}

	compressed := serializedPK
	if algo == dcrec.STEcdsaSecp256k1 || algo == dcrec.STSchnorrSecp256k1 {
		pub, err := secp256k1.ParsePubKey(serializedPK)
		if err != nil {
			return ""
		}
		pubSerComp := pub.SerializeCompressed()

		// Set the y-bit if needed.
		if pubSerComp[0] == 0x03 {
			pubKeyBytes[0] |= (1 << 7)
		}

		compressed = pubSerComp[1:]
	}

	pubKeyBytes = append(pubKeyBytes, compressed...)
	return base58.CheckEncode(pubKeyBytes, netID)
}

type Address interface {
	String() string
	EncodeAddress() string
	ScriptAddress() []byte
	Hash160() *[ripemd160.Size]byte
	IsForNet(*chaincfg.Params) bool
	DSA(*chaincfg.Params) dcrec.SignatureType
	Net() *chaincfg.Params
}

func NewAddressPubKey(decoded []byte, net *chaincfg.Params) (Address, error) {
	if len(decoded) == 33 {
		suite := decoded[0]
		suite &= ^uint8(1 << 7)
		ybit := !(decoded[0]&(1<<7) == 0)
		toAppend := uint8(0x02)
		if ybit {
			toAppend = 0x03
		}

		switch dcrec.SignatureType(suite) {
		case dcrec.STEcdsaSecp256k1:
			return NewAddressSecpPubKey(
				append([]byte{toAppend}, decoded[1:]...),
				net)
		case dcrec.STEd25519:
			return NewAddressEdwardsPubKey(decoded, net)
		case dcrec.STSchnorrSecp256k1:
			return NewAddressSecSchnorrPubKey(
				append([]byte{toAppend}, decoded[1:]...),
				net)
		}
		return nil, ErrUnknownAddressType
	}
	return nil, ErrUnknownAddressType
}

func DecodeAddress(addr string) (Address, error) {
	decoded, netID, err := base58.CheckDecode(addr)
	if err != nil {
		if err == base58.ErrChecksum {
			return nil, ErrChecksumMismatch
		}
		return nil, fmt.Errorf("decoded address is of unknown format: %v",
			err.Error())
	}

	net, err := detectNetworkForAddress(addr)
	if err != nil {
		return nil, ErrUnknownAddressType
	}

	switch netID {
	case net.PubKeyAddrID:
		return NewAddressPubKey(decoded, net)

	case net.PubKeyHashAddrID:
		return NewAddressPubKeyHash(decoded, net, dcrec.STEcdsaSecp256k1)

	case net.PKHEdwardsAddrID:
		return NewAddressPubKeyHash(decoded, net, dcrec.STEd25519)

	case net.PKHSchnorrAddrID:
		return NewAddressPubKeyHash(decoded, net, dcrec.STSchnorrSecp256k1)

	case net.ScriptHashAddrID:
		return NewAddressScriptHashFromHash(decoded, net)

	default:
		return nil, ErrUnknownAddressType
	}
}

func detectNetworkForAddress(addr string) (*chaincfg.Params, error) {
	if len(addr) < 1 {
		return nil, fmt.Errorf("empty string given for network detection")
	}

	return chaincfg.ParamsByNetAddrPrefix(addr[0:1])
}

type AddressPubKeyHash struct {
	net   *chaincfg.Params
	hash  [ripemd160.Size]byte
	netID [2]byte
}

func NewAddressPubKeyHash(pkHash []byte, net *chaincfg.Params, algo dcrec.SignatureType) (*AddressPubKeyHash, error) {
	var addrID [2]byte
	switch algo {
	case dcrec.STEcdsaSecp256k1:
		addrID = net.PubKeyHashAddrID
	case dcrec.STEd25519:
		addrID = net.PKHEdwardsAddrID
	case dcrec.STSchnorrSecp256k1:
		addrID = net.PKHSchnorrAddrID
	default:
		return nil, errors.New("unknown ECDSA algorithm")
	}
	apkh, err := newAddressPubKeyHash(pkHash, addrID)
	if err != nil {
		return nil, err
	}
	apkh.net = net
	return apkh, nil
}

func newAddressPubKeyHash(pkHash []byte, netID [2]byte) (*AddressPubKeyHash, error) {
	if len(pkHash) != ripemd160.Size {
		return nil, errors.New("pkHash must be 20 bytes")
	}
	addr := &AddressPubKeyHash{netID: netID}
	copy(addr.hash[:], pkHash)
	return addr, nil
}

func (a *AddressPubKeyHash) EncodeAddress() string {
	return encodeAddress(a.hash[:], a.netID)
}

func (a *AddressPubKeyHash) ScriptAddress() []byte {
	return a.hash[:]
}

func (a *AddressPubKeyHash) IsForNet(net *chaincfg.Params) bool {
	return a.netID == net.PubKeyHashAddrID ||
		a.netID == net.PKHEdwardsAddrID ||
		a.netID == net.PKHSchnorrAddrID
}

func (a *AddressPubKeyHash) String() string {
	return a.EncodeAddress()
}

func (a *AddressPubKeyHash) Hash160() *[ripemd160.Size]byte {
	return &a.hash
}

func (a *AddressPubKeyHash) DSA(net *chaincfg.Params) dcrec.SignatureType {
	switch a.netID {
	case net.PubKeyHashAddrID:
		return dcrec.STEcdsaSecp256k1
	case net.PKHEdwardsAddrID:
		return dcrec.STEd25519
	case net.PKHSchnorrAddrID:
		return dcrec.STSchnorrSecp256k1
	}
	return -1
}

func (a *AddressPubKeyHash) Net() *chaincfg.Params {
	return a.net
}

type AddressScriptHash struct {
	net   *chaincfg.Params
	hash  [ripemd160.Size]byte
	netID [2]byte
}

func NewAddressScriptHash(serializedScript []byte,
	net *chaincfg.Params) (*AddressScriptHash, error) {
	scriptHash := Hash160(serializedScript)
	ash, err := newAddressScriptHashFromHash(scriptHash, net.ScriptHashAddrID)
	if err != nil {
		return nil, err
	}
	ash.net = net

	return ash, nil
}

func NewAddressScriptHashFromHash(scriptHash []byte,
	net *chaincfg.Params) (*AddressScriptHash, error) {
	ash, err := newAddressScriptHashFromHash(scriptHash, net.ScriptHashAddrID)
	if err != nil {
		return nil, err
	}
	ash.net = net

	return ash, nil
}

func newAddressScriptHashFromHash(scriptHash []byte,
	netID [2]byte) (*AddressScriptHash, error) {
	if len(scriptHash) != ripemd160.Size {
		return nil, errors.New("scriptHash must be 20 bytes")
	}

	addr := &AddressScriptHash{netID: netID}
	copy(addr.hash[:], scriptHash)
	return addr, nil
}

func (a *AddressScriptHash) EncodeAddress() string {
	return encodeAddress(a.hash[:], a.netID)
}

func (a *AddressScriptHash) ScriptAddress() []byte {
	return a.hash[:]
}

func (a *AddressScriptHash) IsForNet(net *chaincfg.Params) bool {
	return a.netID == net.ScriptHashAddrID
}

func (a *AddressScriptHash) String() string {
	return a.EncodeAddress()
}

func (a *AddressScriptHash) Hash160() *[ripemd160.Size]byte {
	return &a.hash
}

func (a *AddressScriptHash) DSA(net *chaincfg.Params) dcrec.SignatureType {
	return -1
}

func (a *AddressScriptHash) Net() *chaincfg.Params {
	return a.net
}

type PubKeyFormat int

const (
	PKFUncompressed PubKeyFormat = iota

	PKFCompressed
)

var ErrInvalidPubKeyFormat = errors.New("invalid pubkey format")

type AddressSecpPubKey struct {
	net          *chaincfg.Params
	pubKeyFormat PubKeyFormat
	pubKey       chainec.PublicKey
	pubKeyHashID [2]byte
}

func NewAddressSecpPubKey(serializedPubKey []byte,
	net *chaincfg.Params) (*AddressSecpPubKey, error) {
	pubKey, err := secp256k1.ParsePubKey(serializedPubKey)
	if err != nil {
		return nil, err
	}

	var pkFormat PubKeyFormat
	switch serializedPubKey[0] {
	case 0x02, 0x03:
		pkFormat = PKFCompressed
	case 0x04:
		pkFormat = PKFUncompressed
	default:
		return nil, ErrInvalidPubKeyFormat
	}

	return &AddressSecpPubKey{
		net:          net,
		pubKeyFormat: pkFormat,
		pubKey:       pubKey,
		pubKeyHashID: net.PubKeyHashAddrID,
	}, nil
}

func (a *AddressSecpPubKey) serialize() []byte {
	switch a.pubKeyFormat {
	default:
		fallthrough
	case PKFUncompressed:
		return a.pubKey.SerializeUncompressed()

	case PKFCompressed:
		return a.pubKey.SerializeCompressed()
	}
}

func (a *AddressSecpPubKey) EncodeAddress() string {
	return encodeAddress(Hash160(a.serialize()), a.pubKeyHashID)
}

func (a *AddressSecpPubKey) ScriptAddress() []byte {
	return a.serialize()
}

func (a *AddressSecpPubKey) Hash160() *[ripemd160.Size]byte {
	h160 := Hash160(a.pubKey.SerializeCompressed())
	array := new([ripemd160.Size]byte)
	copy(array[:], h160)

	return array
}

func (a *AddressSecpPubKey) IsForNet(net *chaincfg.Params) bool {
	return a.pubKeyHashID == net.PubKeyHashAddrID
}

func (a *AddressSecpPubKey) String() string {
	return encodePKAddress(a.serialize(), a.net.PubKeyAddrID,
		dcrec.STEcdsaSecp256k1)
}

func (a *AddressSecpPubKey) Format() PubKeyFormat {
	return a.pubKeyFormat
}

func (a *AddressSecpPubKey) AddressPubKeyHash() *AddressPubKeyHash {
	addr := &AddressPubKeyHash{net: a.net, netID: a.pubKeyHashID}
	copy(addr.hash[:], Hash160(a.serialize()))
	return addr
}

func (a *AddressSecpPubKey) PubKey() chainec.PublicKey {
	return a.pubKey
}

func (a *AddressSecpPubKey) DSA(net *chaincfg.Params) dcrec.SignatureType {
	switch a.pubKeyHashID {
	case net.PubKeyHashAddrID:
		return dcrec.STEcdsaSecp256k1
	case net.PKHSchnorrAddrID:
		return dcrec.STSchnorrSecp256k1
	}
	return -1
}

func (a *AddressSecpPubKey) Net() *chaincfg.Params {
	return a.net
}

func NewAddressSecpPubKeyCompressed(pubkey chainec.PublicKey, params *chaincfg.Params) (*AddressSecpPubKey, error) {
	return NewAddressSecpPubKey(pubkey.SerializeCompressed(), params)
}

type AddressEdwardsPubKey struct {
	net          *chaincfg.Params
	pubKey       chainec.PublicKey
	pubKeyHashID [2]byte
}

func NewAddressEdwardsPubKey(serializedPubKey []byte,
	net *chaincfg.Params) (*AddressEdwardsPubKey, error) {
	pubKey, err := edwards.ParsePubKey(edwards.Edwards(), serializedPubKey)
	if err != nil {
		return nil, err
	}

	return &AddressEdwardsPubKey{
		net:          net,
		pubKey:       pubKey,
		pubKeyHashID: net.PKHEdwardsAddrID,
	}, nil
}

func (a *AddressEdwardsPubKey) serialize() []byte {
	return a.pubKey.Serialize()
}

func (a *AddressEdwardsPubKey) EncodeAddress() string {
	return encodeAddress(Hash160(a.serialize()), a.pubKeyHashID)
}

func (a *AddressEdwardsPubKey) ScriptAddress() []byte {
	return a.serialize()
}

func (a *AddressEdwardsPubKey) Hash160() *[ripemd160.Size]byte {
	h160 := Hash160(a.pubKey.Serialize())
	array := new([ripemd160.Size]byte)
	copy(array[:], h160)

	return array
}

func (a *AddressEdwardsPubKey) IsForNet(net *chaincfg.Params) bool {
	return a.pubKeyHashID == net.PKHEdwardsAddrID
}

func (a *AddressEdwardsPubKey) String() string {
	return encodePKAddress(a.serialize(), a.net.PubKeyAddrID,
		dcrec.STEd25519)
}

func (a *AddressEdwardsPubKey) AddressPubKeyHash() *AddressPubKeyHash {
	addr := &AddressPubKeyHash{net: a.net, netID: a.pubKeyHashID}
	copy(addr.hash[:], Hash160(a.serialize()))
	return addr
}

func (a *AddressEdwardsPubKey) PubKey() chainec.PublicKey {
	return a.pubKey
}

func (a *AddressEdwardsPubKey) DSA(net *chaincfg.Params) dcrec.SignatureType {
	return dcrec.STEd25519
}

func (a *AddressEdwardsPubKey) Net() *chaincfg.Params {
	return a.net
}

type AddressSecSchnorrPubKey struct {
	net          *chaincfg.Params
	pubKey       chainec.PublicKey
	pubKeyHashID [2]byte
}

func NewAddressSecSchnorrPubKey(serializedPubKey []byte,
	net *chaincfg.Params) (*AddressSecSchnorrPubKey, error) {
	pubKey, err := schnorr.ParsePubKey(secp256k1.S256(), serializedPubKey)
	if err != nil {
		return nil, err
	}

	return &AddressSecSchnorrPubKey{
		net:          net,
		pubKey:       pubKey,
		pubKeyHashID: net.PKHSchnorrAddrID,
	}, nil
}

func (a *AddressSecSchnorrPubKey) serialize() []byte {
	return a.pubKey.Serialize()
}

func (a *AddressSecSchnorrPubKey) EncodeAddress() string {
	return encodeAddress(Hash160(a.serialize()), a.pubKeyHashID)
}

func (a *AddressSecSchnorrPubKey) ScriptAddress() []byte {
	return a.serialize()
}

func (a *AddressSecSchnorrPubKey) Hash160() *[ripemd160.Size]byte {
	h160 := Hash160(a.pubKey.Serialize())
	array := new([ripemd160.Size]byte)
	copy(array[:], h160)

	return array
}

func (a *AddressSecSchnorrPubKey) IsForNet(net *chaincfg.Params) bool {
	return a.pubKeyHashID == net.PubKeyHashAddrID
}

func (a *AddressSecSchnorrPubKey) String() string {
	return encodePKAddress(a.serialize(), a.net.PubKeyAddrID,
		dcrec.STSchnorrSecp256k1)
}

func (a *AddressSecSchnorrPubKey) AddressPubKeyHash() *AddressPubKeyHash {
	addr := &AddressPubKeyHash{net: a.net, netID: a.pubKeyHashID}
	copy(addr.hash[:], Hash160(a.serialize()))
	return addr
}

func (a *AddressSecSchnorrPubKey) DSA(net *chaincfg.Params) dcrec.SignatureType {
	return dcrec.STSchnorrSecp256k1
}

func (a *AddressSecSchnorrPubKey) Net() *chaincfg.Params {
	return a.net
}
