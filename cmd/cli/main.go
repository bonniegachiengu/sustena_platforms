package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
)

var bc = blockchain.NewBlockchain()

func main() {
	http.HandleFunc("/create_account", createAccountHandler)
	http.HandleFunc("/transfer", transferHandler)
	http.HandleFunc("/get_balance", getBalanceHandler)
	http.HandleFunc("/get_chain", getChainHandler)
	http.HandleFunc("/get_accounts", getAccountsHandler)

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	account, err := bc.CreateAccount(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(account)
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	var tx blockchain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bc.Transfer(tx.From, tx.To, tx.Amount, tx.Fee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	balance, err := bc.GetBalance(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(balance)
}

func getChainHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc.GetChain())
}

func getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc.GetAccounts())
}
