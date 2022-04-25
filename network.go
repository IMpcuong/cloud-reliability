package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
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

// getNetwork returns the NetWork definition stored in the `config.json` file.
func getNetwork() Network {
	cfg := getNetworkCfg()
	return cfg.Network
}

// getLocalNode returns the LocalNode's data stored in the `config.json` file.
func getLocalNode() Node {
	cfg := getNetworkCfg()
	return cfg.Network.LocalNode
}

// pullNeighborBC pulls the neighbor blockchain from other node in
// the network and connect it with the local node.
func pullNeighborBC() *BlockChain {
	var bc *BlockChain
	Info.Printf("Pulling blockchain from other node in Network...")
	nw := getNetwork()

	for i := 0; i < MAX_ASK_TIME; i++ {
		for _, node := range nw.NeighborNodes {
			if bc == nil || bc.IsEmpty() {
				bc = reqConnectBC(node, nil)
				if bc != nil && !bc.IsEmpty() {
					Info.Printf("Pull blockchain succeeded. Current height: %d", bc.GetDepth())
					return bc
				}
			}
		}
	}
	return bc
}

// ReqConnect send the connection request to the given node.
// Connection succeeded if the given node was existed in the network
// and its address is connectable.
func reqConnectBC(node Node, bc *BlockChain) *BlockChain {
	// Checking if the local node is empty or not.
	// If it's empty, depth is equal `0`.
	var localDepth int
	if bc == nil || bc.GetDepth() == 0 {
		bc = new(BlockChain)
		localDepth = 0
	} else {
		localDepth = bc.GetDepth()
	}

	// Checking if the neighbor node's address is reachable or not.
	neighborDepth, err := getDepthNeighbor(node)
	if err != nil {
		return nil
	}

	// Checking if the local node's current block in the lastest position or not.
	for localDepth < neighborDepth {
		msg := CreateMsgReqBlock(localDepth + 1)
		data := msg.Serialize()

		// Checking if the node address/port is reachable or available.
		conn, err := net.Dial("tcp", node.Address)
		if err != nil {
			Error.Printf("%s is not available!\n", node.Address)
			return nil
		}
		defer conn.Close()

		// Copy the msg bytes data to the connected node.
		_, err = io.Copy(conn, bytes.NewReader(data))
		if err != nil {
			Error.Panic(err)
		}

		// Scan the buffer data and convert it to bytes message.
		scanner := bufio.NewScanner(bufio.NewReader(conn))
		scanner.Scan()
		msgAsBytes := scanner.Bytes()

		// Deserialize the bytes message to `*Message` response.
		msgRes := deserializeMsg(msgAsBytes)
		block := deserializeBlock(msgRes.Data)
		bc.Blocks = append(bc.Blocks, block)
		localDepth++
	}
	return bc
}

// fwHashes forwards the new message's hash data to all neighbor nodes.
func fwHashes(bc *BlockChain) {
	nw := getNetwork()
	for _, node := range nw.NeighborNodes {
		msg := CreateMsgFwHash(bc.GetHashes())
		SendMsg(msg, node)
	}
}

// getDepthNeighbor returns the depth of the given node
// that was connected with local node.
func getDepthNeighbor(node Node) (int, error) {
	msg := CreateMsgReqDepth()

	// Checking if the node address/port is reachable or available.
	conn, err := net.Dial("tcp", node.Address)
	if err != nil {
		Error.Printf("%s is not available!\n", node.Address)
		return 0, err
	}
	defer conn.Close()

	// Copy the msg bytes data to the connected node.
	_, err = io.Copy(conn, bytes.NewReader(msg.Serialize()))
	if err != nil {
		Error.Panic(err)
		return 0, err
	}

	// Scan the buffer data and convert it to bytes message.
	scanner := bufio.NewScanner(bufio.NewReader(conn))
	scanner.Scan()
	msgAsBytes := scanner.Bytes()

	// Deserialize the bytes message to `*Message` response.
	msgRes := deserializeMsg(msgAsBytes)
	neighborDepth, err := strconv.Atoi(string(msgRes.Data))
	if err != nil {
		Error.Printf("Error decoding message %s!", err.Error())
	}
	return neighborDepth, nil
}

// SendMsg send new message to the the given node.
func SendMsg(msg *Message, node Node) {
	conn, err := net.Dial("tcp", node.Address)
	if err != nil {
		Error.Printf("%s is not available!\n", node.Address)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(msg.Serialize()))
	if err != nil {
		Error.Panic(err)
		return
	}
}
