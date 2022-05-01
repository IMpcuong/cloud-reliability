package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/google/go-cmp/cmp"
)

// startBCServer turn on the BlockChain network server.
func startBCServer(bc *Blockchain) {
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
func handleReq(conn net.Conn, bc *Blockchain) {
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
	case CReqHeader:
		handleReqHeader(conn, bc, msg)
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
func handleReqFwHash(conn net.Conn, bc *Blockchain, msg *Message) {
	Info.Printf("BlockChain detected modification. Starting synchronize chain...")
	reqConnectBC(msg.Source, bc)
}

// handleReqDepth handles the request asking for the others node's depth (blockchain)
// for the synchronizing in the local node.
// Response with the message of the other node's depth.
func handleReqDepth(conn net.Conn, bc *Blockchain) {
	resMsg := createMsgResDepth(bc.GetDepth())
	conn.Write(resMsg.Serialize())
}

// handleReqBlock handles the request of pulling block after checking the neighbor node's depth.
// Response with the block was missing and sync it into the local node.
func handleReqBlock(conn net.Conn, bc *Blockchain, msg *Message) {
	reqDepth, err := strconv.Atoi(string(msg.Data))
	if err != nil {
		Error.Print(err.Error())
	}
	block := bc.GetBlockByDepth(reqDepth)
	resMsg := createMsgResBlock(block)
	conn.Write(resMsg.Serialize())
}

// handleReqHeader handles the header identical validation block between local and neighbor node.
func handleReqHeader(conn net.Conn, bc *Blockchain, msg *Message) {
	neighborHeader := deserializeHeader(msg.Data)
	localBlock := bc.GetBlockByDepth(neighborHeader.Depth)
	result := cmp.Equal(*neighborHeader, localBlock.Header)
	resMsg := createMsgResHeader(result)
	conn.Write(resMsg.Serialize())
}

// handlePrintChain handles the request of printing the chain's values in string format.
func handlePrintChain(bc *Blockchain) {
	Info.Printf("%v", bc.Stringify())
}

// handleAddBlock handles the request of adding new block to the chain.
func handleAddBlock(conn net.Conn, bc *Blockchain, msg *Message) {
	block := newBlock(string(msg.Data), bc.GetLatestHash(), bc.GetDepth()+1)
	bc.AddBlock(block)
	fwHashes(bc)
}
