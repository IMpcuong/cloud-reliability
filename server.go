package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

// startBCServer turn on the BlockChain network server.
func startBCServer(bc *BlockChain) {
	cfg := getNetworkCfg()
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
		go handleReq(conn, bc)
	}
}

// handleReq handles all cases of incoming message's command from any connected node.
func handleReq(conn net.Conn, bc *BlockChain) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		Error.Println("Error read data: ", err.Error())
		return
	}

	msg := deserializeMsg(buf[:len])
	Info.Printf("Handle command %s request from port: %s\n", msg.Cmd, conn.RemoteAddr())

	switch msg.Cmd {
	case CFwHashList:
		handleReqFwHash(conn, bc, msg)
	case CReqDepth:
		handleReqDepth(conn, bc)
	case CReqBlock:
		handleReqBlock(conn, bc, msg)
	case CPrintChain:
		handlePrintChain(bc)
	case CAddBlock:
		handleAddBlock(conn, bc, msg)
	default:
		Info.Printf("Command message is invalid!\n")
	}

	conn.Close()
}

// handleReqFwHash handles request forwards hashes list to all neighbor node.
func handleReqFwHash(conn net.Conn, bc *BlockChain, msg *Message) {
	Info.Printf("BlockChain detected modification. Starting synchronize chain...")
	reqConnectBC(msg.Source, bc)
}

// handleReqDepth handles the request asking for the others node's depth (blockchain)
// for the synchronizing in the local node.
// Response with the message of the other node's depth.
func handleReqDepth(conn net.Conn, bc *BlockChain) {
	resMsg := createMsgResDepth(bc.GetDepth())
	conn.Write(resMsg.Serialize())
}

// handleReqBlock handles the request of pulling block after checking the neighbor node's depth.
// Response with the block was missing and sync it into the local node.
func handleReqBlock(conn net.Conn, bc *BlockChain, msg *Message) {
	depth, err := strconv.Atoi(string(msg.Data))
	if err != nil {
		Error.Print(err.Error())
	}
	idx := depth - 1
	block := bc.Blocks[idx]
	resMsg := createMsgResBlock(block)
	conn.Write(resMsg.Serialize())
}

// handlePrintChain handles the request of printing the chain's values in string format.
func handlePrintChain(bc *BlockChain) {
	Info.Printf("%v", bc.Stringify())
}

// handleAddBlock handles the request of adding new block to the chain.
func handleAddBlock(conn net.Conn, bc *BlockChain, msg *Message) {
	bc.AddBlock(string(msg.Data))
	fwHashes(bc)
}
