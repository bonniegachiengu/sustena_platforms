package main

import (
	"fmt"
	"log"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	ops := blockchain.NewBlockchainOperations(bc)

	// Create some accounts
	alice, err := ops.CreateNewAccount()
	if err != nil {
		log.Fatalf("Failed to create account for Alice: %v", err)
	}
	bob, err := ops.CreateNewAccount()
	if err != nil {
		log.Fatalf("Failed to create account for Bob: %v", err)
	}

	// Mint initial supply for Alice
	err = ops.MintInitialSupply(alice, 1000)
	if err != nil {
		log.Fatalf("Failed to mint initial supply for Alice: %v", err)
	}

	// Print initial state
	fmt.Println("Initial State:")
	ops.PrintBlockchainState()

	// Perform a transaction
	err = ops.TransferJoules(alice, bob, 50)
	if err != nil {
		log.Fatalf("Failed to transfer Joules: %v", err)
	}

	// Print final state
	fmt.Println("\nFinal State:")
	ops.PrintBlockchainState()

	// Get account balances
	aliceBalance, err := ops.GetAccountBalance(alice)
	if err != nil {
		log.Fatalf("Failed to get Alice's balance: %v", err)
	}
	bobBalance, err := ops.GetAccountBalance(bob)
	if err != nil {
		log.Fatalf("Failed to get Bob's balance: %v", err)
	}

	fmt.Printf("\nAlice's balance: %.6f Joules\n", aliceBalance)
	fmt.Printf("Bob's balance: %.6f Joules\n", bobBalance)
}
