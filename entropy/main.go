package main

import (
	"fmt"
	"log"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
	"github.com/bonniegachiengu/sustena_platforms/entropy/consensus"
)

func main() {
	fmt.Println("Entropy - Sustena Platform's Blockchain Component")

	// Initialize blockchain
	bc := blockchain.NewBlockchain()

	// Initialize Proof of Stake consensus
	pos := consensus.NewProofOfStake()

	// Add some validators
	pos.AddValidator("Validator1", 100)
	pos.AddValidator("Validator2", 200)
	pos.AddValidator("Validator3", 300)

	// Create and add some blocks
	for i := 1; i <= 5; i++ {
		// Create a sample transaction
		tx := blockchain.Transaction{
			From:   fmt.Sprintf("User%d", i),
			To:     fmt.Sprintf("User%d", i+1),
			Amount: int64(i * 10),
		}

		// Select a validator
		validator := pos.SelectValidator()
		stake := pos.Validators[validator]

		// Create a new block
		lastBlock := bc.GetLastBlock()
		newBlock := blockchain.NewBlock(
			lastBlock.Index+1,
			[]blockchain.Transaction{tx},
			lastBlock.Hash,
			validator,
			stake,
		)

		// Add the block to the blockchain
		err := bc.AddBlock(newBlock)
		if err != nil {
			log.Printf("Error adding block: %v", err)
			continue
		}

		fmt.Printf("Block %d added by validator %s with stake %d\n", i, newBlock.Validator, newBlock.Stake)

		// Validate the block
		if pos.Validate(newBlock) {
			fmt.Println("Block is valid")
		} else {
			fmt.Println("Block is invalid")
		}

		fmt.Println("--------------------")
	}

	// Print the final state of the blockchain
	fmt.Println("Final Blockchain State:")
	for _, block := range bc.Chain {
		fmt.Printf("Block %d: Hash: %s, Validator: %s, Stake: %d\n", 
			block.Index, block.Hash[:10], block.Validator, block.Stake)
	}
}
