package main

import (
	"encoding/json"
	"os"
	"strconv"
)

const (
	CFwHashList = "FW_HASH_LIST" // Forward the source hash list to other nodes.
	CReqDepth   = "REQ_DEPTH"    // Request to fetch the depth of the current node.
	CReqBlock   = "REQ_BLOCK"    // Request to fetch the given block contents.
	CPrintChain = "PRINT_CHAIN"  // Request to print the blockchain from the given node.
	CAddBlock   = "ADD_BLOCK"    // Request to add a new block to the given chain.
	CResDepth   = "RES_DEPTH"    // Response to the requested fetch depth.
	CResBlock   = "RES_BLOCK"    // Response to the requested fetch block contents.
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
		"PRINT_CHAIN",
		"ADD_BLOCK",
		"RES_DEPTH",
		"RES_BLOCK",
	}[pos]
	return codeAsStr
}

// CreateMsgFwHash used to forwards list of hashes from local node to other nodes.
func CreateMsgFwHash(hashes [][]byte) *Message {
	data, err := json.Marshal(hashes)
	if err != nil {
		Error.Panic("Marshal Failed!\n")
		os.Exit(1)
	}
	return &Message{CFwHashList, data, GetLocalNode()}
}

// CreateMsgReqDepth returns a new request message to fetch
// the current depth of a blockchain.
func CreateMsgReqDepth() *Message {
	return &Message{CReqDepth, nil, GetLocalNode()}
}

// CreateMsgReqBlock returns a new request message to fetch a block's contents.
func CreateMsgReqBlock(pos int) *Message {
	return &Message{CReqBlock, []byte(strconv.Itoa(pos)), GetLocalNode()}
}

// CreateMsgResDepth returns a message to response the fetch depth request.
func CreateMsgResDepth(depth int) *Message {
	return &Message{CResDepth, []byte(strconv.Itoa(depth)), GetLocalNode()}
}

// CreateMsgResBlock returns a message to response the fetch block request.
func CreateMsgResBlock(block *Block) *Message {
	return &Message{CResBlock, block.Serialize(), GetLocalNode()}
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

// DeserializeMsg decode the given message from JSON formatter
// into the original data type using `json.Unmarshal()`.
func DeserializeMsg(encoded []byte) *Message {
	msg := new(Message)
	err := json.Unmarshal(encoded, msg)
	if err != nil {
		Error.Printf("Unmarshal message failed!\n")
		os.Exit(1)
	}
	return msg
}
