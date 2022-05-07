package main

import (
	"bytes"
	"math/big"
)

/*
	Base256-to-Base58 conversion (treat both quantities like big-endian).
*/

// NOTE: this alphabet did not include:
// `0` (zero), `O` (`o` uppercase), `I` (`i` uppercase), `l` (`l` lowercase).
var alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// base58Encode returns a byte slice in base58 encoding.
func base58Encode(input []byte) []byte {
	var result []byte

	// x := the `big.Int` of a big-endian bytes input.
	x := big.NewInt(0).SetBytes(input)

	base58 := big.NewInt(int64(len(alphabet)))
	zero := big.NewInt(0)
	// mod := modulus of the arithmetric.
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base58, mod)
		result = append(result, alphabet[mod.Int64()])
	}

	reverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{alphabet[0]}, result...)
		} else {
			break
		}
	}
	return result
}

// base58Decode decodes base58-encoded data.
func base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}
