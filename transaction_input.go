package main

import "fmt"

// Continue the story from the `transaction.go`:
// Basic structure for a TransactionInput.
type TxInput struct {
	TxID      []byte `json:"TxID"`      // TransactionID of the previous qualified transaction.
	TxOutIdx  int    `json:"TxOutIdx"`  // Indexing how many times the buyer has already transferred money.
	Signature []byte `json:"Signature"` // Digital Signature of buyer.
	PubKey    []byte `json:"PubKey"`    // Still the same as the buyer's `PubKeyHash` in `TxOutput`.
}

// Utility functions start from here.

func (txInput *TxInput) Stringify() string {
	txStr := fmt.Sprintf("TxID : %x\n", txInput.TxID)
	txStr += fmt.Sprintf("	+ TxOutIdx   : %d\n", txInput.TxOutIdx)
	txStr += fmt.Sprintf("	+ Signature  : %d\n", txInput.Signature)
	txStr += fmt.Sprintf("	+ Public Key : %d\n", txInput.PubKey)
	return txStr
}
