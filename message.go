package main

import (
	"encoding/json"
	"os"
	"strconv"
)

const (
	CFwHashList = "FW_HASH_LIST" // Request to forward the source hash list to other nodes.
	CReqDepth   = "REQ_DEPTH"    // Request to fetch the depth of the current node.
	CReqBlock   = "REQ_BLOCK"    // Request to fetch the given block contents.
	CReqHeader  = "REQ_HEADER"   // Request to fetch the given block's header to validate against the hash list.
	CPrintChain = "PRINT_CHAIN"  // Request to print the blockchain from the given node.
	CAddBlock   = "ADD_BLOCK"    // Request to add a new block to the given chain.

	CResDepth  = "RES_DEPTH"  // Response to the requested fetch depth.
	CResBlock  = "RES_BLOCK"  // Response to the requested fetch block contents.
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

// Stringify invoked when message stored as enum type.
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

// Response Messages:

// createMsgResDepth returns a message to response the fetch depth request.
func createMsgResDepth(depth int) *Message {
	return createMsg(CResDepth, Itobytes(depth))
}

// createMsgResBlock returns a message to response the fetch block request.
func createMsgResBlock(block *Block) *Message {
	return createMsg(CResBlock, block.Serialize())
}

// createMsgResHeader returns a message containing the result of the checking validation header request.
func createMsgResHeader(isValid bool) *Message {
	return createMsg(CResHeader, []byte(strconv.FormatBool(isValid)))
}

// Utility functions start from here.

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
