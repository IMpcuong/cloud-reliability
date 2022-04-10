package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	// Instantiate new block chain.
	instantiatedChain *BlockChain
	// Allow any action performed once.
	once sync.Once
)

// Simple structure of the blockchain.
type BlockChain struct {
	Blocks []*Block // List of all blocks.
}

// Initialize the first block in the chain.
func InitBlockChain() *BlockChain {
	once.Do(func() {
		instantiatedChain = &BlockChain{[]*Block{NewGenesisBlock("Genesis Block")}}
	})
	return instantiatedChain
}

// Utility functions start from here.

// Checking the chain is empty or not.
func (bc *BlockChain) IsEmpty() bool {
	return bc.Blocks == nil || bc.GetDepth() == 0
}

// Get max length/depth of the chain.
func (bc *BlockChain) GetDepth() int {
	return len(bc.Blocks)
}

// Adding new block from other node to the local chain
// by append local chain's slice with this block.
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[bc.GetDepth()-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// Get the list of all hashes in the blockchain.
func (bc *BlockChain) GetHashes() [][]byte {
	var hashes [][]byte
	for _, block := range bc.Blocks {
		hashes = append(hashes, block.Hash)
	}
	return hashes
}

// Stringify returns a string representation of the chain's values.
func (bc *BlockChain) Stringify() string {
	var chainAsStr string
	for index, block := range bc.Blocks {
		blockAsStr := fmt.Sprintf("%v", block)
		// Convert index number to string with decimal base
		chainAsStr += "[" + strconv.Itoa(index) + "]"
		chainAsStr += blockAsStr
		chainAsStr += "\n"
	}
	return chainAsStr
}

// Seriallize encode the chain's values into JSON formatter using `json.Marshal()`.
func (bc BlockChain) Serialize() []byte {
	encoded, err := json.Marshal(bc)
	if err != nil {
		Error.Printf("Marshal chain failed!\n")
		os.Exit(1)
	}
	return encoded
}

// DeserializeChain decode the chain's values from JSON formatter
// into the original data type using `json.Unmarshal()`.
func DeserializeChain(encoded []byte) *BlockChain {
	bc := new(BlockChain)
	err := json.Unmarshal(encoded, bc)
	if err != nil {
		Error.Printf("Unmarshal chain failed!\n")
		os.Exit(1)
	}
	return bc
}
