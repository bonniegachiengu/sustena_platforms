package main

import (
	"sustena_platforms/entropy/blockchain"
	"sustena_platforms/entropy/cli"
	"sustena_platforms/entropy/mycelium"
	"sustena_platforms/entropy/node"
	"sustena_platforms/utils"
)

const stateFile = "blockchain_state.json"

func main() {
	utils.LogInfo("Starting Sustena Platforms")
	
	// Create P2P network
	network := mycelium.NewP2PNetwork()

	// Create a single node for now
	node := node.NewNode("node1", 8001, stateFile)

	// Add node to the network
	network.AddPeer(node.ID, node.Blockchain)

	// Set P2P network for the node
	node.SetP2PNetwork(network)

	// Add wallet manager to the node
	node.SetWalletManager(blockchain.NewWalletManager())

	// Start the node
	node.Start()

	// Create and run the CLI
	cli := cli.NewCLI(node)
	cli.Run()

	utils.LogInfo("Sustena Platforms shutting down")
}