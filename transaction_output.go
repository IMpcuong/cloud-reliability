package main

// Continue the story from the `transaction.go`:
// Basic structure for a TransactionOutput.
type TxOutput struct {
	// The total amount of currencies that remain intact by the owner before the transaction happens (= 20 Bitcoins).
	Value int

	// The hash value of the public key from buyer A (A owning `Value` before a transaction really happens).
	PubKeyHash []byte
}


