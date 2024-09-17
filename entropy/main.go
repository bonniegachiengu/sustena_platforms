package main

import (
	"fmt"
	"sustena_platforms/entropy/blockchain"
	"sustena_platforms/entropy/mycelium"
	"sustena_platforms/entropy/node"
	"sustena_platforms/utils"
)

const stateFile = "blockchain_state.json"

func main() {
	utils.LogInfo("Starting Sustena Platforms")
	
	// Create P2P network
	network := mycelium.NewP2PNetwork()

	// Create nodes
	node1 := node.NewNode("node1", 8001, stateFile)
	node2 := node.NewNode("node2", 8002, stateFile)
	node3 := node.NewNode("node3", 8003, stateFile)

	// Add nodes to the network
	network.AddPeer(node1.ID, node1.Blockchain)
	network.AddPeer(node2.ID, node2.Blockchain)
	network.AddPeer(node3.ID, node3.Blockchain)

	// Set P2P network for each node
	node1.SetP2PNetwork(network)
	node2.SetP2PNetwork(network)
	node3.SetP2PNetwork(network)

	// Add wallet manager to nodes
	node1.SetWalletManager(blockchain.NewWalletManager())
	node2.SetWalletManager(blockchain.NewWalletManager())
	node3.SetWalletManager(blockchain.NewWalletManager())

	// Start nodes
	node1.Start()
	node2.Start()
	node3.Start()

	// Simulate some transactions and block creation
	simulateTransactions(node1, node2, node3)

	utils.LogInfo("Simulation complete")
}

func simulateTransactions(nodes ...*node.Node) {
	// Create wallets for each node
	for _, n := range nodes {
		wallet := n.WalletManager.CreateWallet()
		address := wallet.GetAddress()
		n.Blockchain.AddValidator(address, 100) // Add validator with 100 stake
		n.Blockchain.Wallets[address] = 1000    // Give each wallet 1000 initial balance
		utils.LogInfo(fmt.Sprintf("Created wallet for node %s: %s with balance 1000 JUL", n.ID, address))
	}

	// Simulate some transactions
	for i := 0; i < 5; i++ {
		sender := nodes[i%len(nodes)]
		receiver := nodes[(i+1)%len(nodes)]
		
		senderWallets := sender.WalletManager.Wallets
		if len(senderWallets) == 0 {
			utils.LogError(utils.NewError(fmt.Sprintf("No wallets found for node %s", sender.ID)))
			continue
		}
		
		// Get the first wallet from the sender's WalletManager
		var senderWallet *blockchain.Wallet
		var senderAddress string
		for address, wallet := range senderWallets {
			senderWallet = wallet
			senderAddress = address
			break
		}
		
		receiverWallets := receiver.WalletManager.Wallets
		if len(receiverWallets) == 0 {
			utils.LogError(utils.NewError(fmt.Sprintf("No wallets found for node %s", receiver.ID)))
			continue
		}
		
		// Get the first wallet address from the receiver's WalletManager
		var receiverAddress string
		for address := range receiverWallets {
			receiverAddress = address
			break
		}

		tx, err := blockchain.NewTransaction(senderWallet, receiverAddress, 10)
		if err != nil {
			utils.LogError(utils.NewError(fmt.Sprintf("Error creating transaction: %v", err)))
			continue
		}

		if sender.Blockchain.AddTransaction(tx) {
			sender.P2P.BroadcastTransaction(sender.ID, tx)
			utils.LogInfo(fmt.Sprintf("Transaction added and broadcasted: %s -> %s, Amount: 10 JUL", senderAddress, receiverAddress))
		} else {
			utils.LogError(utils.NewError(fmt.Sprintf("Failed to add transaction: %s -> %s, Amount: 10 JUL", senderAddress, receiverAddress)))
		}

		// Create a new block
		validator := sender.Blockchain.PoS.SelectValidator()
		if validator != nil {
			err := sender.Blockchain.AddBlock(validator.Address)
			if err != nil {
				utils.LogError(utils.NewError(fmt.Sprintf("Error creating new block: %v", err)))
			} else {
				newBlock := sender.Blockchain.GetLatestBlock()
				sender.PropagateNewBlock(newBlock)
				utils.LogInfo(fmt.Sprintf("New block created and propagated by validator %s", validator.Address))
			}
		} else {
			utils.LogError(utils.NewError("No validator selected for block creation"))
		}

		// Sync blockchains
		for _, n := range nodes {
			n.P2P.SyncBlockchain(n.ID)
		}
	}
}

// Rest of your main.go file...