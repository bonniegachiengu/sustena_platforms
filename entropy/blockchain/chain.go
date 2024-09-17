package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sustena_platforms/entropy/consensus"
	"sustena_platforms/utils"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Blockchain struct {
	Blocks    []*Block
	PoS       *consensus.PoS
	Mempool   []*Transaction
	Wallets   map[string]*Wallet
	db        *sql.DB
}

func NewBlockchain(dbPath string) (*Blockchain, error) {
	// Ensure the directory exists
	err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	bc := &Blockchain{
		PoS:     consensus.NewPoS(),
		Mempool: []*Transaction{},
		Wallets: make(map[string]*Wallet),
		db:      db,
	}

	err = bc.initDB()
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %v", err)
	}

	err = bc.loadBlocksFromDB()
	if err != nil {
		return nil, fmt.Errorf("error loading blocks from database: %v", err)
	}

	return bc, nil
}

func (bc *Blockchain) initDB() error {
	_, err := bc.db.Exec(`
        CREATE TABLE IF NOT EXISTS blocks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ` + "`index`" + ` INTEGER,
            timestamp INTEGER,
            previous_hash TEXT,
            hash TEXT,
            validator TEXT
        );
        CREATE TABLE IF NOT EXISTS transactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            block_id INTEGER,
            tx_id TEXT,
            from_address TEXT,
            to_address TEXT,
            amount REAL,
            fee REAL,
            FOREIGN KEY (block_id) REFERENCES blocks(id)
        );
        CREATE TABLE IF NOT EXISTS wallets (
            address TEXT PRIMARY KEY,
            balance REAL,
            staked REAL
        );
    `)
	return err
}

func (bc *Blockchain) loadBlocksFromDB() error {
	rows, err := bc.db.Query("SELECT id, `index`, timestamp, previous_hash, hash, validator FROM blocks ORDER BY `index`")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var block Block
		err := rows.Scan(&block.ID, &block.Index, &block.Timestamp, &block.PreviousHash, &block.Hash, &block.Validator)
		if err != nil {
			return err
		}
		block.Transactions, err = bc.getTransactionsForBlock(block.ID)
		if err != nil {
			return err
		}
		bc.Blocks = append(bc.Blocks, &block)
	}

	return rows.Err()
}

func (bc *Blockchain) getTransactionsForBlock(blockID int64) ([]*Transaction, error) {
	rows, err := bc.db.Query("SELECT tx_id, from_address, to_address, amount, fee FROM transactions WHERE block_id = ?", blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(&tx.ID, &tx.From, &tx.To, &tx.Amount, &tx.Fee)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

// Add this new method
func (bc *Blockchain) AddValidator(address string, stake float64) {
	if wallet, exists := bc.Wallets[address]; exists {
		bc.PoS.AddValidator(address, stake, wallet.GetTotalBalance())
	}
}

// Modify AddBlock to save to database
func (bc *Blockchain) AddBlock(validator string) error {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	transactions := bc.getTransactionsFromMempool()
	utils.LogInfo(fmt.Sprintf("Creating new block with %d transactions", len(transactions)))
	
	newBlock := CreateBlock(int64(len(bc.Blocks)), transactions, prevBlock.Hash, validator, prevBlock)
	
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
	
	tx, err := bc.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO blocks (`index`, timestamp, previous_hash, hash, validator) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newBlock.Index, newBlock.Timestamp, newBlock.PreviousHash, newBlock.Hash, newBlock.Validator)
	if err != nil {
		tx.Rollback()
		return err
	}

	blockID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, transaction := range newBlock.Transactions {
		_, err = tx.Exec("INSERT INTO transactions (block_id, tx_id, from_address, to_address, amount, fee) VALUES (?, ?, ?, ?, ?, ?)",
			blockID, transaction.ID, transaction.From, transaction.To, transaction.Amount, transaction.Fee)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
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
	err := bc.SaveState()
	if err != nil {
		utils.LogError(fmt.Errorf("failed to save state after adding transaction: %v", err))
	}
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
		if wallet.GetBalance() >= tx.Amount + tx.Fee { // Check if balance covers amount + fee
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
			utils.LogError(utils.NewError(fmt.Sprintf("Insufficient balance for %s: has %.2f, needs %.2f (including fee)", tx.From, wallet.GetBalance(), tx.Amount + tx.Fee)))
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

		fromWallet.DeductBalance(tx.Amount + tx.Fee)
		toWallet.AddBalance(tx.Amount)
		
		// Add the fee to the validator's balance
		validatorWallet.AddBalance(tx.Fee)
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
			bc.PoS.AddValidator(block.Validator, wallet.GetStakedAmount(), wallet.GetTotalBalance())
		}
	}

	return bc, nil
}

// Add this method
func (bc *Blockchain) RegisterWallet(wallet *Wallet) {
	address := wallet.GetAddress()
	bc.Wallets[address] = wallet
	bc.PoS.AddValidator(address, 0, wallet.GetTotalBalance()) // Initialize with 0 stake
	utils.LogInfo(fmt.Sprintf("Registered wallet with address %s", address))
}

// Add this method to the Blockchain struct
func (bc *Blockchain) SaveState() error {
	state := BlockchainState{
		Blocks:  bc.Blocks,
		Wallets: make(map[string]float64),
	}
	for address, wallet := range bc.Wallets {
		state.Wallets[address] = wallet.GetBalance()
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling blockchain state: %v", err)
	}
	return ioutil.WriteFile("blockchain_state.json", data, 0644)
}

func (bc *Blockchain) GetAllBlocks() ([]*Block, error) {
	rows, err := bc.db.Query("SELECT id, `index`, timestamp, previous_hash, hash, validator FROM blocks ORDER BY `index` DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*Block
	for rows.Next() {
		var block Block
		err := rows.Scan(&block.ID, &block.Index, &block.Timestamp, &block.PreviousHash, &block.Hash, &block.Validator)
		if err != nil {
			return nil, err
		}
		block.Transactions, err = bc.getTransactionsForBlock(block.ID)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &block)
	}

	return blocks, rows.Err()
}
