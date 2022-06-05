package main

// A quickly explanation of transmitting mechanism or transaction procedure in a Blockchain system:
//
// 	Transaction: is a transfer of data values that is broadcast to the network and collected into blocks.
//	A transaction typically references previous transaction outputs as new transaction inputs
// 	and dedicates all data values to new outputs.
//
//	We can imagine there is a transaction that has been committed between person A and person B.
//	A typically used 20 Bitcoins to buy an item from B.
//	A has an ID to distinguish him/herself from the others, neither do B.

// The basic structure for a qualified transaction.
type Transaction struct {
	ID     []byte     // Bytes slice to identify the transaction ID itself.
	TxIns  []TxInput  // TransactionInputs array.
	TxOuts []TxOutput // TransactionOutputs array.
}
