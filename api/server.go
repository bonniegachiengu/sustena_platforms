package api

import (
	"encoding/json"
	"net/http"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
)

type Server struct {
	bc  *blockchain.Blockchain
	ops *blockchain.BlockchainOperations
}

type APIConfig struct {
	Blockchain *blockchain.Blockchain
	Operations *blockchain.BlockchainOperations
}

func NewServer(config APIConfig) *Server {
	return &Server{
		bc:  config.Blockchain,
		ops: config.Operations,
	}
}

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	address, err := s.ops.CreateNewAccount(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"address": address})
}

func (s *Server) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
		Fee    int64   `json:"fee"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := s.ops.TransferJoules(req.From, req.To, req.Amount, req.Fee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	balance, err := s.ops.GetAccountBalance(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func (s *Server) GetChainHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.bc.GetChain())
}

func (s *Server) GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.bc.GetAccounts())
}
