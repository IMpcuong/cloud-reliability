package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Simple structure of a block.
type Block struct {
	Hash          []byte // Hash value of each block.
	PrevBlockHash []byte // Previous block's hash value.
	Data          []byte // Data inside the block.
	Timestamp     int64  // Timestamp created the block.
}

// Create Genesis Block (starting point).
func newGenesisBlock(starting string) *Block {
	return newBlock(starting, []byte{})
}

// Create new block for the blockchain.
func newBlock(data string, prevBlockHash []byte) *Block {
	nblock := &Block{[]byte{}, prevBlockHash, []byte(data), time.Now().Unix()}
	nblock.GenHash()
	return nblock
}

// Utility functions start from here.

// Stringify returns a string representation for the given block.
func (block *Block) Stringify() string {
	var blockAsStr string
	blockAsStr += fmt.Sprintf("Previous hash value: %x\n", block.PrevBlockHash)
	blockAsStr += fmt.Sprintf("Data of given block: %x\n", block.Data)
	blockAsStr += fmt.Sprintf("Hash of given block: %x\n", block.Hash)
	return blockAsStr
}

// Create a hash generator for the new block.
func (block *Block) GenHash() {
	bTimeStamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	blockAsBytes := bytes.Join([][]byte{block.PrevBlockHash, block.Data, bTimeStamp}, []byte{})
	hashValue := sha256.Sum256(blockAsBytes)

	block.Hash = hashValue[:]
}

// Serialize encode the given block's value into JSON formatter using `json.Marshal()`.
func (block *Block) Serialize() []byte {
	encoded, err := json.Marshal(block)
	if err != nil {
		Error.Printf("Marshal block failed!\n")
		os.Exit(1)
	}
	return encoded
}

// deserializeBlock decode the given block's value from JSON formatter
// into the original data type using `json.Unmarshal()`.
func deserializeBlock(encoded []byte) *Block {
	block := new(Block)
	err := json.Unmarshal(encoded, block)
	if err != nil {
		Error.Printf("Unmarshal block failed!\n")
		os.Exit(1)
	}
	return block
}
