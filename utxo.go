package main

import (
	"bytes"
	"encoding/hex"

	"github.com/boltdb/bolt"
)

// Un-spend Transaction Output Set - UTxO (The set of remaining transactions output)

const (
	UTXO_BUCKET = "chain_state"
)

type UTxOSet struct {
	Blockchain *Blockchain
}

// Utility functions start from here.

// getBoltProps returns a tuple containing the bucket storage of transactions
// and the cursor indicates the position of a transaction in the bucket.
func getBoltProps(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, *bolt.Cursor) {
	bucket := tx.Bucket(bucketName)
	cursor := bucket.Cursor()

	return bucket, cursor
}

// UTxOSet methods:
// GetUTxOProps returns a tuple containing the transaction's database storage file
// and the bucket's name of its.
func (s UTxOSet) GetUTxOProps() (*bolt.DB, []byte) {
	unspentDB := s.Blockchain.DB
	bucketName := []byte(UTXO_BUCKET)

	return unspentDB, bucketName
}

// FindByPubKey provides a method to search for all of the TxOutputs
// that have associated with the given public hash key (stored in one block).
func (s *UTxOSet) FindByPubKey(pubKeyHash []byte) TxOutputMap {
	db, bucketName := s.GetUTxOProps()
	uTxOs := make(TxOutputMap)

	err := db.View(func(tx *bolt.Tx) error {
		_, cursor := getBoltProps(tx, bucketName)

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txOuts := deserializeTxOutMap(v)

			for idx, txOut := range txOuts {
				if txOut.IsLockedWith(pubKeyHash) {
					uTxOs[idx] = txOut
				}
			}
		}

		return nil
	})
	if err != nil {
		Error.Panic(err)
	}

	return uTxOs
}

// FindSpendableTxOut accesses the bucket storage to retrieve the total amount
// of spendable values from a wallet, also returning the remaining
// TxOutput that can be fulfilled from this wallet too.
func (s UTxOSet) FindSpendableTxOut(pubKeyHash []byte, totalVal int) (int, map[string]TxOutputMap) {
	db, bucketName := s.GetUTxOProps()
	remainTxOuts := make(map[string]TxOutputMap)
	spendableVal := 0

	err := db.View(func(tx *bolt.Tx) error {
		_, cursor := getBoltProps(tx, bucketName)

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txID := hex.EncodeToString(k)
			txOuts := deserializeTxOutMap(v)

			for idxOut, txOut := range txOuts {
				if txOut.IsLockedWith(pubKeyHash) && spendableVal < totalVal {
					spendableVal += txOut.Value
					if remainTxOuts[txID] == nil {
						remainTxOuts[txID] = make(TxOutputMap)
					}
					remainTxOuts[txID][idxOut] = txOut
				}
			}
		}

		return nil
	})
	if err != nil {
		Error.Panic(err)
	}

	return spendableVal, remainTxOuts
}

// CountTxs counting the total number of transactions that were executed.
func (s *UTxOSet) CountTxs() int {
	db, bucketName := s.GetUTxOProps()
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		_, cursor := getBoltProps(tx, bucketName)

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		Error.Panic(err)
	}

	return counter
}

// Rearrange is a helper function that will rebuild a totally new un-spent
// transaction output set with the arrangement as same as the original one.
func (s UTxOSet) Rearrange() {
	db, bucketName := s.GetUTxOProps()

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			Error.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			Error.Panic(err)
		}

		return nil
	})
	if err != nil {
		Error.Panic(err)
	}

	uTxOs := s.Blockchain.FindExistUTxO()
	db.Update(func(tx *bolt.Tx) error {
		bucket, _ := getBoltProps(tx, bucketName)

		for txID, outs := range uTxOs {
			key, err := hex.DecodeString(txID)
			if err != nil {
				Error.Panic(err)
			}

			err = bucket.Put(key, outs.Serialize())
			if err != nil {
				Error.Panic(err)
			}
		}

		return nil
	})
}

// GetAllAddrs returns a list of all addresses information collected from each TxOutput.
func (s UTxOSet) GetAllAddrs() map[string]int {
	db, bucketName := s.GetUTxOProps()
	uTxOs := make(map[string]TxOutputMap)
	addrsInfos := make(map[string]int)

	db.View(func(tx *bolt.Tx) error {
		_, cursor := getBoltProps(tx, bucketName)

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txID := hex.EncodeToString(k)
			txOuts := deserializeTxOutMap(v)
			uTxOs[txID] = txOuts
		}

		return nil
	})

	for _, txOuts := range uTxOs {
		for _, txOut := range txOuts {
			addr := hex.EncodeToString(txOut.PubKeyHash)
			addrsInfos[addr] = txOut.Value
		}
	}
	return addrsInfos
}

// GetTotalValOwnedBy returns the total amount of values owned by a user
// having the provided public key hash.
func (s UTxOSet) GetTotalValOwnedBy(pubKeyHash []byte) int {
	db, bucketName := s.GetUTxOProps()
	totalVal := 0

	db.View(func(tx *bolt.Tx) error {
		_, cursor := getBoltProps(tx, bucketName)

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txOuts := deserializeTxOutMap(v)
			for _, txOut := range txOuts {
				if bytes.Equal(txOut.PubKeyHash, pubKeyHash) {
					totalVal += txOut.Value
				}
			}
		}

		return nil
	})

	return totalVal
}

// VerifyTxIns verifying the integrity of the TxInput set.
func (s UTxOSet) VerifyTxIns(txIns []TxInput) bool {
	db, bucketName := s.GetUTxOProps()
	isValid := true

	db.View(func(tx *bolt.Tx) error {
		bucket, _ := getBoltProps(tx, bucketName)

		for _, txIn := range txIns {
			bytesTxOuts := bucket.Get(txIn.TxID)
			if bytesTxOuts != nil {
				listTxOuts := deserializeTxOutMap(bytesTxOuts)
				if _, ok := listTxOuts[txIn.TxOutIdx]; !ok {
					isValid = false
					return nil
				}
			} else {
				isValid = true
				return nil
			}
		}

		return nil
	})

	return isValid
}

// Update is a helper function born to update the state of all transactions
// that have been stored in the provided block.
func (s UTxOSet) Update(block *Block) {
	db, bucketName := s.GetUTxOProps()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, _ := getBoltProps(tx, bucketName)

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, txIn := range tx.TxIns {
					updatedTxOuts := make(TxOutputMap)
					bytesTxOuts := bucket.Get(txIn.TxID)
					listTxOuts := deserializeTxOutMap(bytesTxOuts)

					for idx, txOut := range listTxOuts {
						if idx != txIn.TxOutIdx {
							updatedTxOuts[idx] = txOut
						}
					}

					if len(updatedTxOuts) == 0 {
						err := bucket.Delete(txIn.TxID)
						if err != nil {
							Error.Panic(err)
						}
					} else {
						err := bucket.Put(txIn.TxID, updatedTxOuts.Serialize())
						if err != nil {
							Error.Panic(err)
						}
					}
				}
			}

			newTxOuts := make(TxOutputMap)
			for idx, txOut := range newTxOuts {
				newTxOuts[idx] = txOut
			}

			err := bucket.Put(tx.ID, newTxOuts.Serialize())
			if err != nil {
				Error.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		Error.Panic(err)
	}
}
