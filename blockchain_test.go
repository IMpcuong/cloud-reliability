package main

import (
	"fmt"
	"testing"
)

func TestInitBC(t *testing.T) {
	var bc *BlockChain = initBlockChain()
	if bc == nil {
		t.Errorf("Cannot initialize block chain!")
	}

	var blocks []*Block = bc.Blocks
	bc.AddBlock("IMpossible send 1 eth/btc to Batman")
	bc.AddBlock("Batman send 2 eth/btc to IMpossible")
	bc.AddBlock("One Punch Man send 3 eth/btc to IMpossible")
	if bc == nil {
		t.Errorf("Cannot add new block to the chain!")
	}
	fmt.Println()

	for _, block := range blocks {
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Timestamp : %x\n", block.Timestamp)
		fmt.Printf("Previous Hash : %x\n", block.PrevBlockHash)
		fmt.Println()
	}
}
