package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sustena_platforms/entropy/consensus"
	"sustena_platforms/utils"
)

type Blockchain struct {
	Blocks    []*Block
	PoS       *consensus.PoS
	Mempool   []*Transaction
	Wallets   map[string]*Wallet // Change this line
}

func NewBlockchain(stateFile string) *Blockchain {
	if _, err := os.Stat(stateFile); err == nil {
		// State file exists, load the blockchain
		bc, err := LoadBlockchainFromDisk(stateFile)
		if err == nil {
			fmt.Println("Loaded existing blockchain state from disk")
			return bc
		}
		fmt.Printf("Error loading blockchain state: %v. Creating new blockchain.\n", err)
	}

	// Create a new blockchain if loading fails or file doesn't exist
	pos := consensus.NewPoS()
	genesisBlock := CreateBlock(0, []*Transaction{}, "", "")
	return &Blockchain{
		Blocks:  []*Block{&genesisBlock},
		PoS:     pos,
		Mempool: []*Transaction{},
		Wallets: make(map[string]*Wallet), // Change this line
	}
}

// Add this new method
func (bc *Blockchain) AddValidator(address string, stake uint64) {
	bc.PoS.AddValidator(address, stake)
}

func (bc *Blockchain) AddBlock(validator string) error {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	transactions := bc.getTransactionsFromMempool()
	utils.LogInfo(fmt.Sprintf("Creating new block with %d transactions", len(transactions)))
	
	newBlock := CreateBlock(int64(len(bc.Blocks)), transactions, prevBlock.Hash, validator)
	
	if !newBlock.IsValid(prevBlock) {
		err := utils.NewError(fmt.Sprintf("Invalid block: failed validation. New block index: %d, Previous block index: %d", newBlock.Index, prevBlock.Index))
		utils.LogError(err)
		return err
	}
	
	if !newBlock.HasValidTransactions() {
		err := utils.NewError("Invalid block: contains invalid transactions")
		utils.LogError(err)
		return err
	}
	
	bc.Blocks = append(bc.Blocks, &newBlock)
	bc.processTransactions(newBlock)

	// Reward the validator
	if validatorWallet, exists := bc.Wallets[validator]; exists {
		validatorWallet.AddBalance(BlockReward)
	}

	utils.LogInfo(fmt.Sprintf("New block added: Index %d, Validator %s, Transactions %d", 
               newBlock.Index, newBlock.Validator, len(newBlock.Transactions)))
	return nil
}

func (bc *Blockchain) AddTransaction(tx *Transaction) bool {
	if !bc.validateTransaction(tx) {
		return false
	}
	
	// Check if transaction is already in mempool
	for _, mempoolTx := range bc.Mempool {
		if mempoolTx.ID == tx.ID {
			utils.LogError(utils.NewError(fmt.Sprintf("Transaction %s is already in the mempool", tx.ID)))
			return false
		}
	}
	
	// Check if transaction is already in a block
	for _, block := range bc.Blocks {
		for _, blockTx := range block.Transactions {
			if blockTx.ID == tx.ID {
				utils.LogError(utils.NewError(fmt.Sprintf("Transaction %s is already in a block", tx.ID)))
				return false
			}
		}
	}
	
	bc.Mempool = append(bc.Mempool, tx)
	return true
}

func (bc *Blockchain) getTransactionsFromMempool() []*Transaction {
	if len(bc.Mempool) > 100 {
		transactions := bc.Mempool[:100]
		bc.Mempool = bc.Mempool[100:]
		return transactions
	}
	transactions := bc.Mempool
	bc.Mempool = []*Transaction{}
	return transactions
}

func (bc *Blockchain) validateTransaction(tx *Transaction) bool {
	if wallet, exists := bc.Wallets[tx.From]; exists {
		if wallet.GetBalance() >= tx.Amount {
			// Convert PublicKeyJSON to ecdsa.PublicKey
			publicKey := &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     tx.PublicKey.X,
				Y:     tx.PublicKey.Y,
			}
			if !VerifyTransaction(tx, publicKey) {
				utils.LogError(utils.NewError(fmt.Sprintf("Transaction %s failed signature verification", tx.ID)))
				return false
			}
			
			// Check for double-spending
			for _, mempoolTx := range bc.Mempool {
				if mempoolTx.ID == tx.ID {
					utils.LogError(utils.NewError(fmt.Sprintf("Transaction %s is already in the mempool", tx.ID)))
					return false
				}
			}
			
			for _, block := range bc.Blocks {
				for _, blockTx := range block.Transactions {
					if blockTx.ID == tx.ID {
						utils.LogError(utils.NewError(fmt.Sprintf("Transaction %s is already in a block", tx.ID)))
						return false
					}
				}
			}
			
			return true
		} else {
			utils.LogError(utils.NewError(fmt.Sprintf("Insufficient balance for %s: has %.2f, needs %.2f", tx.From, wallet.GetBalance(), tx.Amount)))
			return false
		}
	} else {
		utils.LogError(utils.NewError(fmt.Sprintf("Account not found: %s", tx.From)))
		return false
	}
}

func (bc *Blockchain) processTransactions(block Block) {
	for _, tx := range block.Transactions {
		fromWallet := bc.EnsureWallet(tx.From)
		toWallet := bc.EnsureWallet(tx.To)
		validatorWallet := bc.EnsureWallet(block.Validator)

		fromWallet.DeductBalance(tx.Amount)
		toWallet.AddBalance(tx.Amount)
		
		// Add a small fee to the validator as an incentive
		fee := tx.Amount * 0.001 // 0.1% fee
		validatorWallet.AddBalance(fee)
	}
}

func (bc *Blockchain) EnsureWallet(address string) *Wallet {
	if wallet, exists := bc.Wallets[address]; exists {
		return wallet
	}
	newWallet := NewWallet()
	bc.Wallets[address] = newWallet
	return newWallet
}

func (bc *Blockchain) GetBalance(address string) float64 {
	if wallet, exists := bc.Wallets[address]; exists {
		return wallet.GetBalance()
	}
	return 0
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]
		
		if !currentBlock.IsValid(previousBlock) {
			return false
		}
		
		if !currentBlock.HasValidTransactions() {
			return false
		}
	}
	return true
}

func (bc *Blockchain) ValidateChain() error {
	if !bc.IsValid() {
		return fmt.Errorf("blockchain is invalid")
	}
	return nil
}

// Add these new methods

func (bc *Blockchain) ReplaceChain(newChain []*Block) error {
	if len(newChain) <= len(bc.Blocks) {
		return fmt.Errorf("new chain is not longer than the current chain")
	}

	if !bc.IsValidChain(newChain) {
		return fmt.Errorf("new chain is not valid")
	}

	bc.Blocks = newChain
	bc.ReprocessTransactions()
	return nil
}

func (bc *Blockchain) IsValidChain(chain []*Block) bool {
	if len(chain) == 0 {
		return false
	}

	// Check if the first block is the genesis block
	if CalculateHash(*chain[0]) != CalculateHash(*bc.Blocks[0]) {
		return false
	}

	for i := 1; i < len(chain); i++ {
		if !chain[i].IsValid(chain[i-1]) {
			return false
		}
		if !chain[i].HasValidTransactions() {
			return false
		}
	}

	return true
}

func (bc *Blockchain) ReprocessTransactions() {
	// Reset all wallet balances
	for _, wallet := range bc.Wallets {
		wallet.SetBalance(0)
	}

	// Reprocess all transactions in the new chain
	for _, block := range bc.Blocks {
		bc.processTransactions(*block)
	}
}

func (bc *Blockchain) GetChainLength() int {
	return len(bc.Blocks)
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// Add these new methods and structures

type BlockchainState struct {
	Blocks  []*Block
	Wallets map[string]float64
}

func (bc *Blockchain) SaveToDisk(filename string) error {
	walletBalances := make(map[string]float64)
	for address, wallet := range bc.Wallets {
		walletBalances[address] = wallet.GetBalance()
	}
	
	state := BlockchainState{
		Blocks:  bc.Blocks,
		Wallets: walletBalances,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling blockchain state: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing blockchain state to file: %v", err)
	}

	return nil
}

func LoadBlockchainFromDisk(filename string) (*Blockchain, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading blockchain state from file: %v", err)
	}

	var state BlockchainState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling blockchain state: %v", err)
	}

	bc := &Blockchain{
		Blocks:  state.Blocks,
		PoS:     consensus.NewPoS(),
		Mempool: []*Transaction{},
		Wallets: make(map[string]*Wallet),
	}

	// Convert float64 balances to Wallet objects
	for address, balance := range state.Wallets {
		wallet := NewWallet()
		wallet.SetBalance(balance)
		bc.Wallets[address] = wallet
	}

	// Reconstruct the PoS state
	for _, block := range bc.Blocks {
		if wallet, exists := bc.Wallets[block.Validator]; exists {
			bc.PoS.AddValidator(block.Validator, uint64(wallet.GetBalance()))
		}
	}

	return bc, nil
}

// Add this method
func (bc *Blockchain) RegisterWallet(wallet *Wallet) {
	address := wallet.GetAddress()
	bc.Wallets[address] = wallet
	bc.PoS.AddValidator(address, 100) // Add as validator with 100 stake
	utils.LogInfo(fmt.Sprintf("Registered wallet with address %s", address))
}
