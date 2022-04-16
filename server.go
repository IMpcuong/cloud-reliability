package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

// StartBCServer turn on the BlockChain network server.
func StartBCServer(bc *BlockChain) {
	cfg := GetNetworkCfg()
	listener, err := net.Listen("tcp", cfg.Network.LocalNode.Address)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	Info.Println("Local Node listening on port: " + cfg.Network.LocalNode.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			Error.Println("Error accept connection: ", err.Error())
		}
		go HandleReq(conn, bc)
	}
}

// HandleReq handles all cases of incoming message's command from any connected node.
func HandleReq(conn net.Conn, bc *BlockChain) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		Error.Println("Error read data: ", err.Error())
		return
	}

	msg := DeserializeMsg(buf[:len])
	Info.Printf("Handle command %s request from port: %s\n", msg.Cmd, conn.RemoteAddr())

	switch msg.Cmd {
	case CFwHashList:
		HandleReqFwHash(conn, bc, msg)
	case CReqDepth:
		HandleReqDepth(conn, bc)
	case CReqBlock:
		HandleReqBlock(conn, bc, msg)
	case CPrintChain:
		HandlePrintChain(bc)
	case CAddBlock:
		HandleAddBlock(conn, bc, msg)
	default:
		Info.Printf("Command message is invalid!\n")
	}

	conn.Close()
}

// HandleReqFwHash handles request forwards hashes list to all neighbor node.
func HandleReqFwHash(conn net.Conn, bc *BlockChain, msg *Message) {
	Info.Printf("BlockChain detected modification. Starting synchronize chain...")
	ReqConnectBC(msg.Source, bc)
}

// HandleReqDepth handles the request asking for the others node's depth (blockchain)
// for the synchronizing in the local node.
// Response with the message of the other node's depth'.
func HandleReqDepth(conn net.Conn, bc *BlockChain) {
	resMsg := CreateMsgResDepth(bc.GetDepth())
	conn.Write(resMsg.Serialize())
}

// HandleReqBlock handles the request of pulling block after checking the neighbor node's depth.
// Response with the block was missing and sync it into the local node.
func HandleReqBlock(conn net.Conn, bc *BlockChain, msg *Message) {
	depth, err := strconv.Atoi(string(msg.Data))
	if err != nil {
		Error.Print(err.Error())
	}
	idx := depth - 1024
	block := bc.Blocks[idx]
	resMsg := CreateMsgResBlock(block)
	conn.Write(resMsg.Serialize())
}

// HandlePrintChain handles the request of printing the chain's values in string format.
func HandlePrintChain(bc *BlockChain) {
	Info.Printf("%v", bc.Stringify())
}

// HandleAddBlock handles the request of adding new block to the chain.
func HandleAddBlock(conn net.Conn, bc *BlockChain, msg *Message) {
	bc.AddBlock(string(msg.Data))
	FwHashes(bc)
}
