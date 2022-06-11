package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

// NewCoinBaseTx creates a new coin-base transaction. The coin-base transaction can be
// understood as the first transaction that was added in the first block of the chain.
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

// HashTx hashing the whole transaction into a 32 bytes array.
func (tx *Transaction) HashTx() []byte {
	var hash [32]byte

	clonedTx := *tx
	clonedTx.ID = []byte{}
	hash = sha256.Sum256(clonedTx.Serialize())

	return hash[:]
}

// Clone is the method that allows creating a new imitation/emulation
// transaction from the original one.
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

// Sign is a utility function that was invented for the main purpose
// is to help the buyer sign his/her own private key into the exchange deal.
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsCoinbase() {
		return
	}

	clonedTx := tx.Clone()
	signData := fmt.Sprintf("%x", clonedTx)
	// NOTE: Not yet fully understood!
	// IDEA: a full signature was generated from a `rand` number, a buyer's `privKey`,
	// and the corresponding data.
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

// VerifySignature is a helper function that used to verify the
// reliability of a transaction's signature.
func (tx *Transaction) VerifySignature() bool {
	clonedTx := tx.Clone()
	curve := elliptic.P256()

	for _, valIn := range clonedTx.TxIns {
		pubKey := valIn.PubKey

		// HACK: maybe we should use `big.Int{}` constructor instead?
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
		if !ecdsa.Verify(&rawPubKey, []byte(verifyData), r, s) {
			return false
		}
	}
	return true
}

// VerifyValues have similarities in use with the signatures verification method.
// But instead of verifying the signature itself, it checks the balancing between
// the total amount the stream inputs (TxIns) and outputs (TxOuts).
func (tx *Transaction) VerifyValues(prevTxs map[string]Transaction) bool {
	totalIns, totalOuts := 0, 0

	for _, valIn := range tx.TxIns {
		prevTx := prevTxs[hex.EncodeToString(valIn.TxID)]
		totalIns = prevTx.TxOuts[valIn.TxOutIdx].Value
	}

	for _, valOut := range tx.TxOuts {
		totalOuts += valOut.Value
	}

	// Returns true if the total amount of inputs value is equal to
	// the total amount of outputs value.
	return totalIns == totalOuts
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
	strTx := fmt.Sprintf("\n\tID: %x", tx.ID)
	strTx += "\n\tValIn: \n"
	for idx, txIn := range tx.TxIns {
		strTx += fmt.Sprintf("\t[%d]%v\n", idx, txIn)
	}

	strTx += "\n\tValOut: \n"
	for idx, txOut := range tx.TxOuts {
		strTx += fmt.Sprintf("\t[%d]%v\n", idx, txOut)
	}

	return strTx
}
