package main

import (
	"fmt"
	"log"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
)

func main() {
	fmt.Println("Entropy - Sustena Platform's Blockchain Component")

	// Initialize blockchain
	bc := blockchain.NewBlockchain()

	// Create some accounts
	alice, _ := bc.CreateAccount("Alice")
	bob, _ := bc.CreateAccount("Bob")
	charlie, _ := bc.CreateAccount("Charlie")

	// Add some initial balance to Alice's account
	bc.Transfer("GenesisValidator", alice.Address, 1000, 0)

	// Alice stakes some of her balance
	err := bc.AddStake(alice.Address, 500)
	if err != nil {
		log.Printf("Error adding stake: %v", err)
	}

	// Bob stakes some of his balance
	err = bc.Transfer("GenesisValidator", bob.Address, 800, 0)
	if err != nil {
		log.Printf("Error transferring to Bob: %v", err)
	}
	err = bc.AddStake(bob.Address, 300)
	if err != nil {
		log.Printf("Error adding stake for Bob: %v", err)
	}

	// Perform a transaction
	err = bc.Transfer(alice.Address, charlie.Address, 200, 10)
	if err != nil {
		log.Printf("Error transferring: %v", err)
	}

	// Print final balances and stakes
	aliceBalance, _ := bc.GetBalance(alice.Address)
	bobBalance, _ := bc.GetBalance(bob.Address)
	charlieBalance, _ := bc.GetBalance(charlie.Address)

	fmt.Printf("Alice's balance: %d, stake: %d\n", aliceBalance, bc.Accounts[alice.Address].Stake)
	fmt.Printf("Bob's balance: %d, stake: %d\n", bobBalance, bc.Accounts[bob.Address].Stake)
	fmt.Printf("Charlie's balance: %d, stake: %d\n", charlieBalance, bc.Accounts[charlie.Address].Stake)

	// Print the final state of the blockchain
	fmt.Println("Final Blockchain State:")
	for _, block := range bc.Chain {
		fmt.Printf("Block %d: Hash: %s, Validator: %s, Stake: %d, Timestamp: %d\n", 
			block.Index, block.Hash[:10], block.Validator, block.Stake, block.Timestamp)
	}
}
