package main

// Un-spend Transaction Output Set - UTXO (The set of remaining transactions output)

const (
	UTXO_BUCKET = "chain_state"
)

type UTXOSet struct {
	Blockchain *Blockchain
}
