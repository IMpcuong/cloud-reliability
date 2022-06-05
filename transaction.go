package main

// A quickly explanation of transmitting mechanism or transaction procedure in a Blockchain system:
//
// 	Transaction: is a transfer of data values that is broadcast to the network and collected into blocks.
//	A transaction typically references previous transaction outputs as new transaction inputs
// 	and dedicates all data values to new outputs.

// The basic structure for a qualified transaction.
type Transaction struct {
	ID     []byte     // Bytes slice to identify the transaction index itself.
	TxIns  []TxInput  // TransactionInputs array.
	TxOuts []TxOutput // TransactionOutputs array.
}
