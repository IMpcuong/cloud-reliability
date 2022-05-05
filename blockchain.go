package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/vrecan/death"
)

const (
	DB_FILE       = "blockchain.db"
	BLOCKS_BUCKET = "blocks"
)

// Simple structure of the Blockchain.
type Blockchain struct {
	DB *bolt.DB // BoltDB database stored the blockchain.
}

// Iterator implementation for the Blockchain.
type BlockchainIter struct {
	CurHash    []byte      // Hash value of the current block.
	Blockchain *Blockchain // Blockchain itself.
}

// Initialize an empty blockchain and save it to the `DB_FILE`
// if this file is not present yet.
func initBlockChain(node string) *Blockchain {
	absPath := getAbsPathDB(node)
	if dbExist(absPath) {
		fmt.Println("Blockchain database is already exists!")
		return nil
	}

	// Open the database storage file with `read-write` permission.
	db, err := bolt.Open(absPath, 0600, nil)
	if err != nil {
		Error.Fatal(err)
	}

	// `read-write` permission is granted to the user.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
		if err != nil {
			Error.Fatal(err)
		}
		return nil
	})
	if err != nil {
		Error.Fatal(err)
	}

	return &Blockchain{DB: db}
}

// Utility functions start from here.

// Checking the chain is empty or not.
func (bc *Blockchain) IsEmpty() bool {
	return len(bc.GetLatestHash()) == 0
}

// Iterator returns a new iterable blockchain.
func (bc *Blockchain) Iterator() *BlockchainIter {
	lastHash := bc.GetLatestHash()
	bcIter := &BlockchainIter{
		CurHash:    lastHash,
		Blockchain: bc,
	}
	return bcIter
}

// GetDepth returns max height/depth of the blockchain.
func (bc *Blockchain) GetDepth() int {
	var lastBlock *Block

	// Managed the read-only transaction to retrieve the value corresponding with the `l` key.
	err := bc.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET)) // retrieves bucket by its name.
		lastHash := bucket.Get([]byte("l"))        // `l` was defined as key of the latest block's hash.
		if lastHash == nil {
			return nil
		}

		// If the last block's hash value exists, decode its hash value.
		blockData := bucket.Get(lastHash)
		lastBlock = deserializeBlock(blockData)
		return nil
	})
	if err != nil {
		Error.Panic(err)
		return 0
	}

	// Returns depth equal to zero if `lastBlock` does not exist
	// else returns the last block's current position.
	if lastBlock == nil {
		return 0
	}
	return lastBlock.Header.Depth
}

// GetLatestHash retrieves the latest block's hash value.
func (bc *Blockchain) GetLatestHash() []byte {
	// Latest hash value.
	var latest []byte

	// Managed the read-only transaction to retrieve the value corresponding with the `l` key.
	err := bc.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET)) // retrieves bucket by its name.
		latest = bucket.Get([]byte("l"))           // `l` was defined as key of the latest block's hash.
		return nil
	})
	if err != nil {
		Error.Panic(err)
	}

	return latest
}

// GetBlockByDepth returns the block at the given depth/height.
func (bc *Blockchain) GetBlockByDepth(depth int) *Block {
	bcIter := bc.Iterator()
	for {
		curBlock := bcIter.Next()
		if curBlock.Header.Depth == depth {
			return curBlock
		}
		if curBlock.IsGenesis() {
			break
		}
	}
	return nil
}

// Adding a new given block from another node or this local node itself to the local chain
// by appending the local chain's slice with this block.
func (bc *Blockchain) AddBlock(block *Block) {
	pow := newProofOfWork(block)

	if !pow.Validate() {
		nonce, hash := pow.Run()
		block.Header.Nonce = nonce
		block.Header.Hash = hash
	}

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		if bc.IsEmpty() {
			bc.PutBlock(bucket, block.Header.Hash, block.Serialize())
		} else {
			// `l` was defined as key of the latest block's hash.
			lastHash := bucket.Get([]byte("l"))
			// retrieves the encoded data from the last block.
			encodedLastBlock := bucket.Get(lastHash)
			// decodes the last block to retrieves the latest `*Block`.
			lastBlock := deserializeBlock(encodedLastBlock)

			if block.Header.Depth > lastBlock.Header.Depth &&
				bytes.Equal(block.Header.PrevBlockHash, lastBlock.Header.PrevBlockHash) {
				bc.PutBlock(bucket, block.Header.Hash, block.Serialize())
			} else {
				Error.Printf("Block is invalid! Failed to add block: \n%v\n", block)
				Error.Printf("Current latest block: \n%v\n", lastBlock)
			}
		}

		return nil
	})

	if err != nil {
		Error.Panic(err)
	}
}

// PutBlock sets 2 pairs:
// 	`(key, value)` = `(hash, data)`,
// 	`(key, value)` = `("l", latest_hash)` (special pair)
// from the newest block into the bucket (both pairs store inside latest block).
// Remember that bucket is the place where all transactions are stored.
// Each transaction is the pair of a key (block's hash) and a value (block's data)
// except the special pair.
func (bc *Blockchain) PutBlock(bucket *bolt.Bucket, hash, data []byte) {
	err := bucket.Put(hash, data)
	if err != nil {
		Error.Panic(err)
	}
	err = bucket.Put([]byte("l"), hash)
	if err != nil {
		Error.Panic(err)
	}
}

// Get the list of all hashes in the blockchain.
func (bc *Blockchain) GetHashes() [][]byte {
	var hashes [][]byte

	bcIter := bc.Iterator()
	// Iterating through all blocks and return the list of hashes until reaching the genesis ones.
	for {
		block := bcIter.Next()
		hashes = append(hashes, block.Header.Hash)
		if block.IsGenesis() {
			break
		}
	}

	return hashes
}

// Stringify returns a string representation of the chain's values.
func (bc *Blockchain) Stringify() string {
	var chainAsStr string

	bcIter := bc.Iterator()
	// Iterating through all blocks and return the blockchain as string representation
	// until reaching the genesis ones.
	for {
		block := bcIter.Next()
		blockAsStr := fmt.Sprintf("%v", block)
		// Convert index number to string with decimal base
		chainAsStr += "[" + strconv.Itoa(block.Header.Depth) + "]"
		chainAsStr += blockAsStr
		chainAsStr += "\n"
		if block.IsGenesis() {
			break
		}
	}
	return chainAsStr
}

// Seriallize encode the chain's values into JSON formatter using `json.Marshal()`.
func (bc Blockchain) Serialize() []byte {
	encoded, err := json.Marshal(bc)
	if err != nil {
		Error.Printf("Marshal chain failed!\n")
		os.Exit(1)
	}
	return encoded
}

// deserializeChain decode the chain's values from JSON formatter
// into the original data type using `json.Unmarshal()`.
func deserializeChain(encoded []byte) *Blockchain {
	bc := new(Blockchain)
	err := json.Unmarshal(encoded, bc)
	if err != nil {
		Error.Printf("Unmarshal chain failed!\n")
		os.Exit(1)
	}
	return bc
}

// Next iterate over the blockchain by each block's hash value
// with reverse order until reaching the genesis block.
func (iter *BlockchainIter) Next() *Block {
	var block *Block

	// Managed the read-only transaction to retrieve the value corresponding with the block's hash (key).
	err := iter.Blockchain.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET)) // retrieves bucket its by name.
		encoded := bucket.Get(iter.CurHash)        // encoded block's data.
		block = deserializeBlock(encoded)          // decoded block's data.
		return nil
	})

	if err != nil {
		Error.Panic(err)
	}

	// Assigning the current hash with the previous hash.
	iter.CurHash = block.Header.PrevBlockHash
	return block
}

// dbExist returns true if the given DB's file path/name was existed.
func dbExist(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// getLocalBC retrieves entire local blockchain information from the `DB_FILE`
// if this file exists.
func getLocalBC(node string) *Blockchain {
	absPath := getAbsPathDB(node)
	if !dbExist(absPath) {
		return nil
	}

	// Open or create a new database storage file with `read-write` permission.
	// NOTE: Bolt cannot access multiple proccesses the same database at the same time.
	db, err := bolt.Open(absPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		Error.Fatal(err)
	}

	return &Blockchain{DB: db}
}

// closeDB forces the database to be closed.
func closeDB(bc *Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		bc.DB.Close()
	})
}

// getAbsPathDB returns the absolute path to the database storage file in the given node.
//@@@ FIXME: this is a temporary solution, maybe automatically later.
func getAbsPathDB(node string) string {
	absPath := filepath.Join("config/", node, "/", DB_FILE)
	return absPath
}
