package main

import (
	"testing"
)

func TestInitBC(t *testing.T) {
	bc := initBlockChain()
	if bc.IsEmpty() {
		t.Errorf("Cannot initialize block chain!")
	}
	// NOTE: import to avoid static check of unused code.
	deserializeChain(bc.Serialize())

	// bc.AddBlock("IMpossible send 1 eth/btc to Batman")
	// bc.AddBlock("Batman send 2 eth/btc to IMpossible")
	// bc.AddBlock("One Punch Man send 3 eth/btc to IMpossible")
	// if bc == nil {
	// 	t.Errorf("Cannot add new block to the chain!")
	// } else {
	// 	chain := deserializeChain(bc.Serialize())
	// 	fmt.Printf("Chain value: %v\n", chain)
	// }
	// fmt.Println()

	// for _, block := range bc.Blocks {
	// 	if block != nil {
	// 		fmt.Printf("Hash : %x\n", block.Hash)
	// 		fmt.Printf("Data : %s\n", block.Data)
	// 		fmt.Printf("Timestamp : %x\n", block.Timestamp)
	// 		fmt.Printf("Previous Hash : %x\n", block.PrevBlockHash)
	// 		fmt.Println()
	// 	}
	// }
}
