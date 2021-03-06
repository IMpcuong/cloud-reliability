package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	CFwHashList = "FW_HASH_LIST" // Request to forward the source hash list to other nodes.
	CReqDepth   = "REQ_DEPTH"    // Request to fetch the depth of the current node.
	CReqBlock   = "REQ_BLOCK"    // Request to fetch the given block contents.
	CReqHeader  = "REQ_HEADER"   // Request to fetch the given block's header to validate against the hash list.
	CReqAddr    = "REQ_ADDR"     // Request to get node's address.
	CReqPrf     = "REQ_PRF"      // Request to get block's proof.
	CPrintChain = "PRINT_CHAIN"  // Request to print the blockchain from the given node.
	CAddBlock   = "ADD_BLOCK"    // Request to add a new block to the given chain.
	CAddTx      = "ADD_TX"       // Request to add a new transaction to the provided block.

	CResDepth  = "RES_DEPTH"  // Response to the requested fetch depth.
	CResBlock  = "RES_BLOCK"  // Response to the requested fetch block contents.
	CResTx     = "RES_Tx"     // Response to the requested adding new transaction to the provided block.
	CResAddr   = "RES_ADDR"   // Response to the requested fetch node's address.
	CResPrf    = "RES_PRF"    // Response to the validate block's proof request.
	CResHeader = "RES_HEADER" // Response to the requested fetch header validation code with block's data.
)

// Using when commands stored as enums type.
type MsgCmd int

// `Message` is the method that's describe how data exchange between each node.
type Message struct {
	Cmd    string `json:"cmd"`      // Request command.
	Data   []byte `json:"data"`     // Contents of message.
	Source Node   `json:"src_node"` // Contents from source node.
}

// Utility functions start from here.

// Stringify invoked when message stored as enum type (unused yet).
func (pos MsgCmd) Stringify() string {
	codeAsStr := [...]string{
		"FW_HASH_LIST",
		"REQ_DEPTH",
		"REQ_BLOCK",
		"REQ_HEADER",
		"PRINT_CHAIN",
		"ADD_BLOCK",
		"RES_DEPTH",
		"RES_BLOCK",
		"RES_HEADER",
	}[pos]
	return codeAsStr
}

// createMsg is the common method for creating a new message with a given code and data.
func createMsg(cmd string, data []byte) *Message {
	return &Message{
		Cmd:    cmd,
		Data:   data,
		Source: getLocalNode(),
	}
}

// Request Messages:

// createMsgFwHash used to forwards list of hashes from local node to other nodes.
func createMsgFwHash(hashes [][]byte) *Message {
	data, err := json.Marshal(hashes)
	if err != nil {
		Error.Panic("Marshal Failed!\n")
		os.Exit(1)
	}
	return createMsg(CFwHashList, data)
}

// createMsgReqDepth returns a new request message to fetch
// the current depth of a blockchain.
func createMsgReqDepth() *Message {
	return createMsg(CReqDepth, []byte{})
}

// createMsgReqBlock returns a new request message to fetch a block's contents
// from a specific position.
func createMsgReqBlock(pos int) *Message {
	return createMsg(CReqBlock, Itobytes(pos))
}

// createMsgReqHeader returns a message to validate the given block's header.
func createMsgReqHeader(header Header) *Message {
	return createMsg(CReqHeader, header.Serialize())
}

// NOTE: unused yet!
// createMsgReqAddr returns a new request message to fetch a node's address.
func CreateMsgReqAddr() *Message {
	return createMsg(CReqAddr, []byte{})
}

// @@@
func createMsgReqPrf(prf []byte) *Message {
	return createMsg(CReqPrf, prf)
}

// createMsgReqAddTx creates a new message to request adding a transaction to a new block.
func createMsgReqAddTx(tx *Transaction) *Message {
	return createMsg(CAddTx, tx.Serialize())
}

// Response Messages:

// createMsgResDepth returns a message to response the fetch depth request.
func createMsgResDepth(depth int) *Message {
	return createMsg(CResDepth, Itobytes(depth))
}

// createMsgResBlock returns a message to response the fetch block request.
func createMsgResBlock(block *Block) *Message {
	return createMsg(CResBlock, block.Serialize())
}

// createMsgResAddTx returns a message to response the adding transaction request.
func createMsgResAddTx(isSuccess bool) *Message {
	return createMsg(CResTx, []byte(strconv.FormatBool(isSuccess)))
}

// createMsgResHeader returns a message containing the result of the checking validation header request.
func createMsgResHeader(isValid bool) *Message {
	return createMsg(CResHeader, []byte(strconv.FormatBool(isValid)))
}

// createMsgResAddr returns a message to response the fetch node's address request.
func createMsgResAddr() *Message {
	return createMsg(CResAddr, []byte{})
}

// @@@
func createMsgResPrf(isValid bool) *Message {
	return createMsg(CResPrf, []byte(strconv.FormatBool(isValid)))
}

// Utility functions start from here.

func (msg *Message) Export(path string) {
	prettyMarshal, e := json.MarshalIndent(msg, "", "  ")
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}

	e = ioutil.WriteFile(path, prettyMarshal, 0644)
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}
}

// Serialize encode the given message into JSON formatter using `json.Marshal()`.
func (msg *Message) Serialize() []byte {
	encoded, err := json.Marshal(msg)
	if err != nil {
		Error.Printf("Marshal message failed!\n")
		os.Exit(1)
	}
	return encoded
}

// deserializeMsg decode the given message from JSON formatter
// into the original data type using `json.Unmarshal()`.
func deserializeMsg(encoded []byte) *Message {
	msg := new(Message)
	err := json.Unmarshal(encoded, msg)
	if err != nil {
		Error.Printf("Unmarshal message failed!\n")
		os.Exit(1)
	}
	return msg
}
