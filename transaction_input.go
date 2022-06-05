package main

// Basic structure for a TransactionInput.
type TxInput struct {
	TxID      []byte
	TxOutID   int
	Signature []byte
	PubKey    []byte
}
