package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const (
	PUB_KEY_PREFIX    = byte(0x04)
	NW_VERSION        = byte(0x00)
	ADDR_CHECKSUM_LEN = 4
)

// Wallet contains a public-private keypair that can be used to identify itself.
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

// WalletJson is used to store the Wallet data structure in the JSON file.
type WalletJson struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

// Utility functions start from here.

// newWallet returns a new Wallet instance.
func newWallet() *Wallet {
	privKey, pubKey := newKeyPair()
	wallet := Wallet{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Address:    genAddr(pubKey),
	}
	return &wallet
}

// newKeyPair generates a new keypair for the wallet structure.
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		Error.Panic(err)
	}
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)

	return *privKey, pubKey
}

/*
Simple imitation schema for generating new `Address` in Bitcoin network (Pk := `PublicKey`)

Schema:
	. nwVersion
	. ripemd160(sha256(Pk)) -> Pk_Hash
	. sha256(sha256(nw_Version + Pk_Hash))[:4] -> checksum
	--------------------------------------------------------------
	base58Encode(nwVersion + Pk_hash + checksum) -> Wallet_Address
*/
func genAddr(pubKey []byte) string {
	version := []byte{NW_VERSION}
	pubKeyHash := hashPubKey(pubKey)
	versionPayload := append(version, pubKeyHash...)
	checksum := checksum(versionPayload)

	// payload := nwVersion + Pk_Hash + checksum
	// Original length of the address hash value.
	payload := append(versionPayload, checksum...)

	// Convert payload from 256 bytes decrease to 58 bytes.
	address := string(base58Encode(payload))
	return address
}

// hashPubKey returns the hash value of the public key by using `ripemd160` hasher.
func hashPubKey(pubKey []byte) []byte {
	pubKeySHA := sha256.Sum256(pubKey)

	// RIPEMD-160 is a hash function computes a 160-bits messages digest.
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubKeySHA[:])
	if err != nil {
		Error.Panic(err)
	}

	pubRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return pubRIPEMD160
}

// checksum returns the checksum of `PublicKey` after hashing through `sha256.Sum256()` twice.
func checksum(hash []byte) []byte {
	firstSHA := sha256.Sum256(hash)
	secondSHA := sha256.Sum256(firstSHA[:])

	// Returns the first 4 characters in the checksum hash value.
	return secondSHA[:ADDR_CHECKSUM_LEN]
}

// validateAddr checks if the wallet address is valid.
func validateAddr(address string) bool {
	payload := base58Decode([]byte(address))
	actualChecksum := payload[len(payload)-ADDR_CHECKSUM_LEN:]

	version := payload[0]
	pubKeyHash := payload[1 : len(payload)-ADDR_CHECKSUM_LEN]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	// Checking the actual checksum versus the expected checksum value.
	return bytes.Equal(actualChecksum, targetChecksum)
}

// Wallet's methods:

// ToJson converts the `Wallet` instance to a JSON storage file.
func (w *Wallet) ToJson() *WalletJson {
	wJson := new(WalletJson)
	wJson.PrivateKey = hex.EncodeToString(w.PrivateKey.D.Bytes())
	wJson.PublicKey = hex.EncodeToString(w.PublicKey)
	wJson.Address = w.Address
	return wJson
}

// Stringify returns a string representation for the provided `Wallet` data.
func (w Wallet) Stringify() string {
	bs, err := json.MarshalIndent(w, "", "   ")
	if err != nil {
		Error.Fatal(err)
	}
	return string(bs)
}

// WalletJson's methods:

// ToWallet converts the JSON payload to a `Wallet` instance.
func (wj *WalletJson) ToWallet() *Wallet {
	w := new(Wallet)
	curve := elliptic.P256()
	privKeyAsBytes, err := hex.DecodeString(wj.PrivateKey)
	if err != nil {
		Error.Fatal(err)
	}

	w.PrivateKey.D = new(big.Int).SetBytes(privKeyAsBytes)
	w.PrivateKey.PublicKey.Curve = curve
	w.PrivateKey.PublicKey.X, w.PrivateKey.PublicKey.Y = curve.ScalarBaseMult(privKeyAsBytes)
	w.PublicKey, err = hex.DecodeString(wj.PublicKey)
	if err != nil {
		Error.Fatal(err)
	}

	w.Address = wj.Address
	return w
}

// Stringify returns a string representation for the given `WalletJson` instance.
func (wj WalletJson) Stringify() string {
	strWallet := "\n  ** Wallet Information ** \n"
	strWallet += fmt.Sprintf("  + Private Key (%d bytes) : %s\n", len(wj.PrivateKey), wj.PrivateKey)
	strWallet += fmt.Sprintf("  + Public Key (%d bytes) : %s\n", len(wj.PublicKey), wj.PublicKey)
	strWallet += fmt.Sprintf("  + Address (%d bytes) : %s\n", len(wj.Address), wj.Address)
	return strWallet
}
