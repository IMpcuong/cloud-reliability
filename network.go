package main

const (
	MAX_ASK_TIME = 2
)

// In the starting point of this project this struct just only contains
// the `Address` of each `Node` itself.
type Node struct {
	Address string `json:"address"`
}

// Define the connection between the local node was running
// and the others node in the P2P network.
type Network struct {
	LocalNode     Node   `json:"local_node"`     // Node itself was running int the chain.
	NeighborNodes []Node `json:"neighbor_nodes"` // Other nodes were connected in the network.
}

// Utility functions start from here.

// GetNetwork returns the NetWork definition stored in the `config.json` file.
func GetNetwork() Network {
	cfg := GetNetworkCfg()
	return cfg.Network
}

// GetLocalNode returns the LocalNode definition stored in the `config.json` file.
func GetLocalNode() Node {
	cfg := GetNetworkCfg()
	return cfg.Network.LocalNode
}

// PullNeighborBC pulls the neighbor blockchain from other node in
// the network and connect it with the local node.
func PullNeighborBC() *BlockChain {
	var bc *BlockChain
	Info.Printf("Pulling blockchain from other node in Network...")
	nw := GetNetwork()

	for i := 0; i < MAX_ASK_TIME; i++ {
		for _, node := range nw.NeighborNodes {
			if bc == nil || bc.IsEmpty() {
				bc = ReqConnect(node, nil)
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
// Connection succeeded if the node given existed in the network.
func ReqConnect(node Node, bc *BlockChain) *BlockChain {
	// var localDepth int
	// if bc == nil || bc.GetDepth() == 0 {
	// 	bc = new(BlockChain)
	// 	localDepth = 0
	// } else {
	// 	localDepth = bc.GetDepth()
	// }

	// neighborDepth, err := GetDepthNeighbor(node)
	return new(BlockChain)
}

func GetDepthNeighbor(node Node) (int, error) {
	return 0, nil
}
