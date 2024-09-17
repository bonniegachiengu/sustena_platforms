package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sustena_platforms/entropy/blockchain"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Blockchain *blockchain.Blockchain
	router     *gin.Engine
}

func NewServer(bc *blockchain.Blockchain) *Server {
	router := gin.Default()
	server := &Server{Blockchain: bc, router: router}
	server.setupRoutes()
	return server
}

func (s *Server) Start(port int) {
	// Serve static files from the React app build directory
	webDir := filepath.Join("..", "web", "build")
	fmt.Printf("Serving static files from: %s\n", webDir)
	s.router.Static("/static", filepath.Join(webDir, "static"))

	// Serve index.html for the root path and any other unmatched routes
	s.router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(webDir, "index.html"))
	})

	// API routes are already set up in setupRoutes()

	fmt.Printf("Server starting on port %d\n", port)
	err := s.router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		api.GET("/blockchain", s.getBlockchain)
		api.GET("/blocks/:index", s.getBlock)
		api.POST("/createWallet", s.createWallet)
		api.GET("/getBalance/:address", s.getBalance)
		api.POST("/sendTransaction", s.sendTransaction)
		api.POST("/forgeBlock", s.forgeBlock)
		api.POST("/stakeJUL", s.stakeJUL)
		api.POST("/unstakeJUL", s.unstakeJUL)
		api.GET("/getCommunityFund", s.getCommunityFund)
	}
}

func (s *Server) getBlockchain(c *gin.Context) {
	c.JSON(http.StatusOK, s.Blockchain.Blocks)
}

func (s *Server) getBlock(c *gin.Context) {
	// Implementation remains the same
}

func (s *Server) createWallet(c *gin.Context) {
	wallet := blockchain.NewWallet()
	s.Blockchain.RegisterWallet(wallet)
	c.JSON(http.StatusOK, gin.H{"address": wallet.GetAddress()})
}

func (s *Server) getBalance(c *gin.Context) {
	address := c.Param("address")
	balance := s.Blockchain.GetBalance(address)
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (s *Server) sendTransaction(c *gin.Context) {
	var req struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	wallet := s.Blockchain.Wallets[req.From]
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender wallet not found"})
		return
	}
	
	tx, err := blockchain.NewTransaction(wallet, req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if s.Blockchain.AddTransaction(tx) {
		c.JSON(http.StatusOK, gin.H{"message": "Transaction added successfully", "txID": tx.ID})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transaction"})
	}
}

func (s *Server) forgeBlock(c *gin.Context) {
	validator := s.Blockchain.PoS.SelectValidator()
	if validator == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No validator available"})
		return
	}
	
	err := s.Blockchain.AddBlock(validator.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	newBlock := s.Blockchain.GetLatestBlock()
	c.JSON(http.StatusOK, gin.H{"message": "New block forged", "block": newBlock})
}

func (s *Server) stakeJUL(c *gin.Context) {
	var req struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	wallet := s.Blockchain.Wallets[req.Address]
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}
	
	if err := wallet.Stake(req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	s.Blockchain.PoS.AddValidator(req.Address, req.Amount, wallet.GetTotalBalance())
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully staked %.2f JUL", req.Amount)})
}

func (s *Server) unstakeJUL(c *gin.Context) {
	var req struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	wallet := s.Blockchain.Wallets[req.Address]
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}
	
	if err := wallet.Unstake(req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	s.Blockchain.PoS.UpdateStake(req.Address, wallet.GetStakedAmount(), wallet.GetTotalBalance())
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully unstaked %.2f JUL", req.Amount)})
}

func (s *Server) getCommunityFund(c *gin.Context) {
	fundBalance := s.Blockchain.PoS.GetCommunityFund()
	c.JSON(http.StatusOK, gin.H{"balance": fundBalance})
}
