package main

type Transaction struct {
	ID     []byte
	TxIns  []TxInput
	TxOuts []TxOutput
}
