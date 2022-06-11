package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	MAX_ASK_TIME = 1
)

// In the starting point of this project this struct just only contains
// the `Address` of each `Node` itself.
type Node struct {
	// Address of the node itself.
	Address string `json:"address"`
}

// Define the connection between the local node was running
// and the others node in the P2P network.
type Network struct {
	// Node itself was running/mentioning in the network.
	LocalNode Node `json:"local_node"`
	// Other nodes were connected in the network.
	NeighborNodes []Node `json:"neighbor_nodes"`
}

// Utility functions start from here.

// getNetwork returns the Network definition stored in the `config.json` file.
func getNetwork() Network {
	cfg := getNetworkCfg()
	return cfg.Network
}

// getLocalNode returns the LocalNode's data stored in the `config.json` file.
func getLocalNode() Node {
	cfg := getNetworkCfg()
	return cfg.Network.LocalNode
}

// syncNeighborBC pulls the neighbor blockchain from the other node in
// the network and synchronizes the local node with it.
func syncNeighborBC(bc *Blockchain) {
	Info.Printf("Pulling blockchain from other node in Network...")
	nw := getNetwork()

	for i := 0; i < MAX_ASK_TIME; i++ {
		for _, node := range nw.NeighborNodes {
			Info.Printf("Try to synchronize with node: %v", node.Address)
			if reqConnectBC(node, bc) {
				Info.Printf("Sync blockchain succeeded. Current height: %d", bc.GetDepth())
				return
			}
		}
	}
}

// reqConnectBC send the connection request to the given node.
// Connection succeeded if the given node existed in the network,
// its address is connectable, and the owner has fewer blocks identical to the other.
func reqConnectBC(node Node, bc *Blockchain) bool {
	// Checking if the local node is empty or not.
	// If it's empty, depth is equal `0`.
	localDepth := bc.GetDepth()

	// Retrieve the neighbor node's depth/length.
	neighborDepth, err := getDepthNeighbor(node)
	if err != nil {
		return false
	}

	Info.Printf("Depth comparison between [local - neighbor]: [%v - %v]", localDepth, neighborDepth)

	detectIdentical(node, bc, localDepth, neighborDepth)
	Info.Println()

	syncNode(node, bc, localDepth, neighborDepth)
	return true
}

func detectIdentical(node Node, bc *Blockchain, local, neighbor int) {
	minDepth := minVal(local, neighbor)

	// Compare the identical minimum of blocks from both sides.
	// NOTE: block position starts from index 1 not 0 like usual case.
	for pos := 1; pos <= minDepth; pos++ {
		if isIdentical := cmpBlockWithNeighbor(bc.GetBlockByDepth(pos), node); isIdentical {
			Info.Printf("Block [%d] similarity detects completed. Progress: %d%%", pos, pos*100/minDepth)
		} else {
			Error.Fatalf("Block [%d] detected distinction. Exit prompt!", pos)
			os.Exit(1)
		}
	}
}

// syncNode pulls/synchronize blocks from any side if the opposite side
// has more blocks than the other.
func syncNode(node Node, bc *Blockchain, local, neighbor int) {
	if local < neighbor {
		Info.Printf("Pull [%d] blocks from neighbor node", neighbor-local)
		for pos := local + 1; pos <= neighbor; pos++ {
			pullBlockNeighbor(bc, node, pos)
			Info.Printf("Pulled block [%d] completed. Progress: %d%%", pos, pos*100/neighbor)
		}
	} else {
		// TODO: implement the pulling process in the case the neighbor node has fewer blocks than the local.
		// NOTE: 1. Need to detect which neighbor node has more blocks than the local.
		//       2. Then figure out the number of blocks needed to pull from the neighbor.
		panic("Not implemented yet!")
	}
}

// cmpBlockWithNeighbor returns true if blocks in the same position have identical data.
func cmpBlockWithNeighbor(block *Block, node Node) bool {
	msg := createMsgReqHeader(block.Header)
	data := msg.Serialize()

	conn, err := openConn(node)
	err = transferData(conn, data)
	if err != nil {
		Error.Panic(err)
		return false
	}

	// Scan the buffer data and convert it to bytes message.
	scannerData, scannable := scanData(conn)
	if !scannable {
		return false
	}

	// Deserialize the bytes message to `*Message` response.
	msgRes := deserializeMsg(scannerData)
	if isValid := strParsable(msgRes.Data); isValid {
		return isValid
	}
	return false
}

// pullBlockNeighbor pulls some limited amount of blocks from the neighbor node
// if their blockchain's total length is bigger than the local.
func pullBlockNeighbor(bc *Blockchain, node Node, posBlock int) {
	msg := createMsgReqBlock(posBlock)
	data := msg.Serialize()

	conn, err := openConn(node)
	err = transferData(conn, data)
	if err != nil {
		Error.Panic(err)
	}

	// Scan the buffer data and convert it to bytes message.
	scannerData, scannable := scanData(conn)
	if !scannable {
		return
	}

	// Deserialize the bytes message to `*Message` response.
	msgRes := deserializeMsg(scannerData)
	block := deserializeBlock(msgRes.Data)

	// Adding new block to the current node's blockchain.
	bc.AddBlock(block)
}

// fwHashes forwards the new message's hash data to all neighbor nodes.
func fwHashes(bc *Blockchain) {
	nw := getNetwork()
	for _, node := range nw.NeighborNodes {
		msg := createMsgFwHash(bc.GetHashes())
		sendMsg(msg, node)
	}
}

// getDepthNeighbor returns the depth of the given node
// that was connected with local node.
func getDepthNeighbor(node Node) (int, error) {
	msg := createMsgReqDepth()
	data := msg.Serialize()

	conn, err := openConn(node)
	err = transferData(conn, data) // @@@ FIXME
	if err != nil {
		Error.Panic(err)
		return 0, err
	}

	scannerData, scannable := scanData(conn)
	if !scannable {
		return 0, nil
	}

	// Deserialize the bytes message to `*Message` response.
	msgRes := deserializeMsg(scannerData)
	neighborDepth := Bytestoi(msgRes.Data)

	return neighborDepth, nil
}

// sendMsg send new message to the the given node.
func sendMsg(msg *Message, node Node) {
	conn, err := openConn(node)
	data := msg.Serialize()

	err = transferData(conn, data)
	if err != nil {
		Error.Panic(err)
		return
	}
}

// checkPort returns true if the connection to the given port was established.
func checkPort(host, port string) bool {
	timeout := time.Duration(3) * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened!", net.JoinHostPort(host, port))
	}
	return true
}

// openConn checking if the node's address/port is reachable or available
// at the moment request connection was established.
func openConn(node Node) (net.Conn, error) {
	conn, err := net.Dial("tcp", node.Address)
	if err != nil {
		Error.Printf("%s is not available!\n", node.Address)
		return nil, err
	}
	defer conn.Close()

	return conn, nil
}

// transferData copy the message bytes data to the provided/connected node.
func transferData(conn net.Conn, data []byte) error {
	readerData := bytes.NewReader(data)
	_, err := io.Copy(conn, readerData)
	if err != nil {
		Error.Panic(err)
		return err
	}

	return nil
}

// scanData scan the buffer data and convert it to bytes message.
func scanData(conn net.Conn) ([]byte, bool) {
	reader := bufio.NewReader(conn)     // Read data from given connection.
	scanner := bufio.NewScanner(reader) // Scan the data has been read from the connection.

	scannerData := scanner.Bytes() // Bytes message.
	scannable := scanner.Scan()

	return scannerData, scannable
}

// strParsable parse the bytes data to string format.
// Then validate the message response header and check if it's valid.
func strParsable(data []byte) bool {
	isValid, err := strconv.ParseBool(string(data))
	if err != nil {
		Error.Printf("Parse failed!")
		return false
	}

	return isValid
}
