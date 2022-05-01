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

type Block struct {
	Header Header `json:"BlockHeader"` // Header contains identity of one block.
	Data   []byte `json:"Data"`        // Data contained by one block.
}

// Simple structure of a block.
type Header struct {
	PrevBlockHash []byte `json:"PrevBlockHash"` // Previous block's hash value.
	Hash          []byte `json:"Hash"`          // Hash value of each block.
	Timestamp     int64  `json:"Timestamp"`     // Timestamp created the block.
	Depth         int    `json:"Depth"`         // Position or current depth of each block.
	Nonce         int    `json:"Nonce"`         // Number only used once.
}

// Create Genesis Block (starting point).
func newGenesisBlock() *Block {
	return newBlock("Genesis Block", []byte{}, 1)
}

// Create and append new block to the chain.
func newBlock(data string, prevBlockHash []byte, curDepth int) *Block {
	nHeader := Header{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte(data),
		Timestamp:     time.Now().Unix(),
		Depth:         curDepth,
		Nonce:         0,
	}
	nblock := &Block{nHeader, []byte{}}
	nblock.GenHash()
	return nblock
}

// Utility functions start from here.

// Block's methods:
// IsGenesis returns true if the block is genesis block.
func (block *Block) IsGenesis() bool {
	return len(block.Header.PrevBlockHash) == 0
}

// Stringify returns a string representation for the given block.
func (block *Block) Stringify() string {
	var blockAsStr string
	blockAsStr += fmt.Sprintf("Previous hash value: %x\n", block.Header.PrevBlockHash)
	blockAsStr += fmt.Sprintf("Block's Data: %x\n", block.Data)
	blockAsStr += fmt.Sprintf("Block's Hash: %x\n", block.Header.Hash)
	blockAsStr += fmt.Sprintf("Block's Depth: %x\n", block.Header.Hash)
	blockAsStr += fmt.Sprintf("Block's Nonce: %x\n", block.Header.Hash)
	blockAsStr += fmt.Sprintf("Timestamp created this block: %x\n", block.Header.Hash)
	return blockAsStr
}

// Create a hash generator for the new block.
func (block *Block) GenHash() {
	bTimeStamp := []byte(strconv.FormatInt(block.Header.Timestamp, 10))
	headerAsBytes := bytes.Join([][]byte{block.Header.PrevBlockHash, block.Data, bTimeStamp}, []byte{})
	headerHashVal := sha256.Sum256(headerAsBytes)

	block.Header.Hash = headerHashVal[:]
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

// Header's methods:
// Serialize encode the given block's header into JSON formatter using `json.Marshal()`.
func (header *Header) Serialize() []byte {
	encoded, err := json.Marshal(header)
	if err != nil {
		Error.Printf("Marshal header failed!\n")
		os.Exit(1)
	}
	return encoded
}

// deserializeHeader decode the given block's header value from JSON formatter
// into the original data type using `json.Unmarshal()`.
func deserializeHeader(encoded []byte) *Header {
	header := new(Header)
	err := json.Unmarshal(encoded, header)
	if err != nil {
		Error.Printf("Unmarshal header failed!\n")
		os.Exit(1)
	}
	return header
}
