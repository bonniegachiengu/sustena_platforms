package main

import (
	"sustena_platforms/entropy/blockchain"
	"sustena_platforms/api"
	"sustena_platforms/utils"
	"fmt"
	"path/filepath"
	"os"
)

func main() {
	utils.LogInfo("Starting Sustena Blockchain application...")

	entropyDir := filepath.Join(".", "entropy")
	err := os.MkdirAll(entropyDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating entropy directory: %v\n", err)
		return
	}

	dbPath := filepath.Join(entropyDir, "blockchain.db")
	fmt.Printf("Initializing blockchain with database path: %s\n", dbPath)

	bc, err := blockchain.NewBlockchain(dbPath)
	if err != nil {
		fmt.Printf("Error creating blockchain: %v\n", err)
		return
	}

	fmt.Println("Blockchain initialized successfully")

	server := api.NewServer(bc)
	fmt.Println("Starting Sustena Blockchain server...")
	server.Start(8080)
}