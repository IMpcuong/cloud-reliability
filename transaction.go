package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
)

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

const (
	// The subsidy value that had been given to new user the first time they were joining the network.
	SUBSIDY = 25
)

// Utility functions start from here.

func NewCoinBaseTx(addrTo string) *Transaction {
	txIn := TxInput{[]byte{}, -1, nil, []byte{}}
	txOut := newTxOut(SUBSIDY, addrTo)
	coinbaseTX := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	coinbaseTX.ID = coinbaseTX.HashTx()

	return &coinbaseTX
}

// IsCoinbase checking if the transaction is the coinbase transaction
// or more specifically if this transaction stored in the first block.
func (tx Transaction) IsCoinbase() bool {
	return len(tx.TxIns) == 1 &&
		len(tx.TxIns[0].TxID) == 0 &&
		tx.TxIns[0].TxOutIdx == -1
}

func (tx *Transaction) HashTx() []byte {
	var hash [32]byte

	clonedTx := *tx
	clonedTx.ID = []byte{}
	hash = sha256.Sum256(clonedTx.Serialize())

	return hash[:]
}

func (tx *Transaction) Clone() Transaction {
	var txIns []TxInput
	var txOuts []TxOutput

	for _, valIn := range tx.TxIns {
		txIns = append(txIns, TxInput{
			TxID:      valIn.TxID,
			TxOutIdx:  valIn.TxOutIdx,
			Signature: nil,
			PubKey:    valIn.PubKey,
		})
	}

	for _, valOut := range tx.TxOuts {
		txOuts = append(txOuts, TxOutput{
			Value:      valOut.Value,
			PubKeyHash: valOut.PubKeyHash,
		})
	}

	clonedTx := Transaction{tx.ID, txIns, txOuts}
	return clonedTx
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsCoinbase() {
		return
	}

	clonedTx := tx.Clone()
	signData := fmt.Sprintf("%x", clonedTx)
	// NOTE: Not yet fully understood!
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(signData))
	if err != nil {
		Error.Panic(err)
	}

	// NOTE: the main point of this process represents cloning the data of a transaction from a block.
	// Executing all of the necessary calculations on the cloned transaction,
	// before returning the signature to the original one.
	signature := append(r.Bytes(), s.Bytes()...)
	for idx := range clonedTx.TxIns {
		clonedTx.TxIns[idx].Signature = nil
		tx.TxIns[idx].Signature = signature
	}
}

func (tx *Transaction) VerifySignature() bool {
	clonedTx := tx.Clone()
	curve := elliptic.P256()

	for _, valIn := range clonedTx.TxIns {
		pubKey := valIn.PubKey

		// Expected values:
		r := new(big.Int)
		s := new(big.Int)
		signLen := len(valIn.Signature)
		r.SetBytes(valIn.Signature[:(signLen / 2)])
		s.SetBytes(valIn.Signature[(signLen / 2):])

		// Real signature values:
		x := new(big.Int)
		y := new(big.Int)
		keyLen := len(pubKey) - 1
		x.SetBytes(pubKey[1:(keyLen/2 + 1)])
		y.SetBytes(pubKey[(keyLen/2 + 1):])

		// Transaction data that needs to be verified.
		verifyData := fmt.Sprintf("%x", clonedTx)

		rawPubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}
		if ecdsa.Verify(&rawPubKey, []byte(verifyData), r, s) == false {
			return false
		}
	}
	return true
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	if err != nil {
		Error.Panic(err)
	}

	return encoded.Bytes()
}

func DeserializeTx(data []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)
	if err != nil {
		Error.Panic(err)
	}

	return &tx
}

func (tx Transaction) Stringify() string {
	strTx := fmt.Sprintf("\n	ID: %x", tx.ID)
	strTx += fmt.Sprintf("\n	ValIn :\n")
	for idx, txIn := range tx.TxIns {
		strTx += fmt.Sprintf("	[%d]%v\n", idx, txIn)
	}

	strTx += fmt.Sprintf("\n	ValOut :\n")
	for idx, txOut := range tx.TxOuts {
		strTx += fmt.Sprintf("	[%d]%v\n", idx, txOut)
	}

	return strTx
}
