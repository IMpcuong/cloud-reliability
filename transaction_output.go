package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Continue the story from the `transaction.go`:
// Basic structure for a TransactionOutput.
type TxOutput struct {
	// The total amount of currencies that remain intact by the owner before the transaction happens (= 20 Bitcoins).
	Value int

	// The hash value of the public key from buyer A (A owning `Value` before a transaction really happens).
	PubKeyHash []byte
}

// Utility functions start from here.

// newTxOut creates a new TxOutput with the provided value and a nil public key hash.
func newTxOut(val int, addr string) *TxOutput {
	nTxOutput := &TxOutput{
		Value:      val,
		PubKeyHash: nil,
	}
	nTxOutput.LockTx(addr)

	return nTxOutput
}

// LockTx depicts the progression of a transaction that is already
// occupied by a buyer and identify by using their unique identifier hash.
func (txOut *TxOutput) LockTx(addr string) {
	decodeAddr := base58Decode([]byte(addr))
	buyerHash := decodeAddr[1 : len(addr)-4]

	// Locking a transaction with the buyer is PubKeyHash.
	txOut.PubKeyHash = buyerHash
}

// IsLocked returns true if the transaction is locked with the buyer's public key hash.
func (txOut *TxOutput) IsLocked(buyerHash []byte) bool {
	return bytes.Equal(txOut.PubKeyHash, buyerHash)
}

func (txOut *TxOutput) Stringify() string {
	str := fmt.Sprintf("Value : %d\n", txOut.Value)
	str += fmt.Sprintf("PubKeyHash : %x ", txOut.PubKeyHash)
	return str
}

// Map of list of all available TxOutput.
type TxOutputMap map[int]TxOutput

func (txOutMap *TxOutputMap) SerializeTxOutMap() []byte {
	var buf bytes.Buffer

	encode := gob.NewEncoder(&buf)
	err := encode.Encode(txOutMap)
	if err != nil {
		Error.Panic(err)
	}

	return buf.Bytes()
}

func DeserializeTxOutMap(data []byte) TxOutputMap {
	var txOutMap TxOutputMap

	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&txOutMap)
	if err != nil {
		Error.Panic(err)
	}

	return txOutMap
}
