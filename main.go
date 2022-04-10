package main

import "fmt"

func main() {
	bc := InitBlockChain()
	msg := new(Message)
	fmt.Println(msg)

	bc.AddBlock("IMpossible send 1 eth/btc to Batman")
	bc.AddBlock("Batman send 2 eth/btc to IMpossible")
	bc.AddBlock("One Punch Man send 3 eth/btc to IMpossible")

	fmt.Println()

	for _, block := range bc.Blocks {
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Timestamp : %x\n", block.Timestamp)
		fmt.Printf("Previous Hash : %x\n", block.PrevBlockHash)
		fmt.Println()
	}
}
