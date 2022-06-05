package main

type TxInput struct {
	TxID      []byte
	TxOutID   int
	Signature []byte
	PubKey    []byte
}
