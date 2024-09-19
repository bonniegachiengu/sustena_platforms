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
	Router     *gin.Engine
}

func NewServer(bc *blockchain.Blockchain) *Server {
	router := gin.Default()
	server := &Server{Blockchain: bc, Router: router}
	server.setupRoutes()
	return server
}

func (s *Server) Start(port int) {
	webDir := filepath.Join("..", "web", "build")
	fmt.Printf("Serving static files from: %s\n", webDir)
	s.Router.Static("/static", filepath.Join(webDir, "static"))
	s.Router.StaticFile("/favicon.ico", filepath.Join(webDir, "favicon.ico"))

	// Serve index.html for the root path and any other unmatched routes
	s.Router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(webDir, "index.html"))
	})

	fmt.Printf("Server starting on port %d\n", port)
	err := s.Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func (s *Server) setupRoutes() {
	api := s.Router.Group("/api")
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
		api.POST("/purchaseJUL", s.purchaseJUL)
		api.GET("/getValidators", s.getValidators)
		api.GET("/mempool", s.getMempool)
	}
}

func (s *Server) getBlockchain(c *gin.Context) {
	blocks, err := s.Blockchain.GetAllBlocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blocks)
}

func (s *Server) getBlock(c *gin.Context) {
	// Implementation remains the same
}

func (s *Server) createWallet(c *gin.Context) {
	wallet := blockchain.NewWallet()
	s.Blockchain.RegisterWallet(wallet)
	
	// Give the new wallet some initial balance (e.g., 100 JUL)
	wallet.AddBalance(100)
	
	c.JSON(http.StatusOK, gin.H{
		"address": wallet.GetAddress(),
		"balance": wallet.GetBalance(),
	})
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

	if wallet.GetBalance() < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}
	
	tx, err := blockchain.NewTransaction(wallet, req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if s.Blockchain.AddTransaction(tx) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction added successfully",
			"txID": tx.ID,
			"fee": tx.Fee,
			"totalAmount": tx.Amount + tx.Fee,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transaction"})
	}
}

func (s *Server) forgeBlock(c *gin.Context) {
	validator := s.Blockchain.PoS.SelectValidator()
	if validator == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No validator available"})
		return
	}
	
	err := s.Blockchain.AddBlock(validator.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	newBlock := s.Blockchain.GetLatestBlock()
	
	// Reward the validator
	validatorWallet := s.Blockchain.Wallets[validator.Address]
	if validatorWallet != nil {
		blockReward := blockchain.GetBlockReward()
		validatorWallet.AddBalance(blockReward)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "New block forged",
		"block": newBlock,
		"validator": validator.Address,
		"reward": blockchain.GetBlockReward(),
	})
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

func (s *Server) purchaseJUL(c *gin.Context) {
	var req struct {
		Address   string  `json:"address"`
		USDAmount float64 `json:"usdAmount"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	wallet, exists := s.Blockchain.Wallets[req.Address]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}
	
	julPurchased := wallet.PurchaseJUL(req.USDAmount)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Purchased %.2f JUL", julPurchased),
		"julAmount": julPurchased,
	})
}

func (s *Server) getValidators(c *gin.Context) {
	validators := s.Blockchain.PoS.GetValidators()
	c.JSON(http.StatusOK, gin.H{"validators": validators})
}

func (s *Server) getMempool(c *gin.Context) {
	c.JSON(http.StatusOK, s.Blockchain.GetMempool())
}
