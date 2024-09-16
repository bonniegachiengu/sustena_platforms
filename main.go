package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
	"github.com/bonniegachiengu/sustena_platforms/api"
)

func main() {
	bc := blockchain.NewBlockchain()
	ops := blockchain.NewBlockchainOperations(bc)

	// Initialize API server
	server := api.NewServer(api.APIConfig{
		Blockchain: bc,
		Operations: ops,
	})

	// Set up routes
	http.HandleFunc("/create_account", server.CreateAccountHandler)
	http.HandleFunc("/transfer", server.TransferHandler)
	http.HandleFunc("/get_balance", server.GetBalanceHandler)
	http.HandleFunc("/get_chain", server.GetChainHandler)
	http.HandleFunc("/get_accounts", server.GetAccountsHandler)
	http.HandleFunc("/buy_jul", server.BuyJULHandler)
	http.HandleFunc("/get_mempool", server.GetMempoolHandler)

	// Serve static files for the UI
	fs := http.FileServer(http.Dir("./blockchain-ui/build"))
	http.Handle("/", fs)

	// Start the server
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
