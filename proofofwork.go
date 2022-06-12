package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const (
	// An arbitrary difficulty number needs to be satisfied to mine a new block.
	DIFFICULTY = 16
	// Maximum value can be reached by a block's nonce number.
	MAX_NONCE = math.MaxInt64
)

// Proof of Work algorithm structure.
type ProofOfWork struct {
	Block  *Block   // Block needed to be validate.
	Target *big.Int // Upper bound of block's hash value.
}

// Initialize the Proof of Work default structure.
func newProofOfWork(block *Block) *ProofOfWork {
	// NOTE: will convert hash to `bigInt` and check if it's less than the target later.
	target := big.NewInt(1)
	// Shift the default target (=1) to the left 240 bits,
	// now the target equal to `2^240` (~30 bytes).
	target.Lsh(target, uint(256-DIFFICULTY))

	pow := &ProofOfWork{
		Block:  block,
		Target: target,
	}
	return pow
}

// Utility functions start from here.

// PrepareData generates the data that will be used to digest by the `SHA256` algorithm.
// This function will be consuming the incremented `nonce` as the argument,
// combining `nonce` with the block's data that we expected to be accomplishing the constraint.
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	txAsBytes := []byte{}
	for _, tx := range pow.Block.Transactions {
		txAsBytes = append(txAsBytes, tx.Serialize()...)
	}

	// Concatenate all the needed data to a bytes slice.
	data := bytes.Join(
		[][]byte{
			pow.Block.Header.PrevBlockHash,
			txAsBytes,
			Itobytes(int(pow.Block.Header.Timestamp)),
			Itobytes(pow.Block.Header.Depth),
			// Nonce is the incremented counter needs to be found.
			Itobytes(nonce),
		},
		[]byte{},
	)
	return data
}

// Run is the execution function or the core of the PoW algorithm.
// This function is used to find the satisfied `nonce` to mine a new block
// with brute force approach and also returns the corresponded hash value.
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte // 32 bytes = 256 bits.
	nonce := 0

	Info.Printf("Mining process starting...")

COUNTER: // for loop of the nonce counter in the range [0, 2^63-1]
	for nonce < MAX_NONCE {
		data := pow.PrepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		// Cmp function in `bigInt` := `hashInt < pow.Target == true`
		if hashInt.Cmp(pow.Target) == -1 {
			break COUNTER
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")
	return nonce, hash[:]
}

// Validate checks the satisfaction of the block's hash value against the target constraint.
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.Block.Header.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	// Returns true if the `hashInt` value is less than the `target` number.
	return hashInt.Cmp(pow.Target) == -1
}
