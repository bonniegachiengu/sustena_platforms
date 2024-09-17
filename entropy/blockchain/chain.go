package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Blockchain struct {
	Chain    []*Block
	Accounts map[string]*Account
	Mempool  []Transaction
	Exchange *Account
	mu       sync.RWMutex
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:    []*Block{GenesisBlock()},
		Accounts: make(map[string]*Account),
		Mempool:  []Transaction{},
	}
	bc.initializeExchange()
	go bc.startBlockTimer()
	return bc
}

func (bc *Blockchain) initializeExchange() {
	exchangeAccount, _ := bc.CreateAccount("Exchange")
	bc.Exchange = exchangeAccount
	bc.MintInitialSupply(bc.Exchange.Address, InitialSupply)
}

func (bc *Blockchain) startBlockTimer() {
	for {
		time.Sleep(BlockTime)
		bc.createNewBlock()
	}
}

func (bc *Blockchain) createNewBlock() {
	if len(bc.Mempool) > 0 {
		lastBlock := bc.Chain[len(bc.Chain)-1]
		validator := bc.selectValidator()
		totalFees := int64(0)
		for _, tx := range bc.Mempool {
			totalFees += tx.Fee
		}

		newBlock := &Block{
			Index:        lastBlock.Index + 1,
			Timestamp:    time.Now().Unix(),
			Transactions: bc.Mempool,
			PrevHash:     lastBlock.Hash,
			Validator:    validator,
			Stake:        bc.Accounts[validator].Stake,
		}
		newBlock.Hash = calculateHash(newBlock.Index, newBlock.Timestamp, newBlock.Transactions, newBlock.PrevHash, newBlock.Validator, newBlock.Stake)
		
		// Apply transactions
		for _, tx := range newBlock.Transactions {
			bc.applyTransaction(tx)
		}

		bc.Chain = append(bc.Chain, newBlock)
		bc.Mempool = []Transaction{}

		// Add validator reward transaction to the mempool for the next block
		rewardTx := Transaction{
			From:   "System",
			To:     validator,
			Amount: totalFees,
			Fee:    0,
		}
		bc.Mempool = append(bc.Mempool, rewardTx)
	}
}

func (bc *Blockchain) applyTransaction(tx Transaction) {
	if tx.From != "System" {
		bc.Accounts[tx.From].Balance -= tx.Amount + tx.Fee
	}
	bc.Accounts[tx.To].Balance += tx.Amount
}

func (bc *Blockchain) Transfer(from, to string, amount, fee int64) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	sender, exists := bc.Accounts[from]
	if !exists {
		return errors.New("sender account not found")
	}

	if sender.Balance < amount+fee {
		return errors.New("insufficient balance")
	}

	tx := Transaction{
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    fee,
	}

	bc.Mempool = append(bc.Mempool, tx)
	return nil
}

func GenesisBlock() *Block {
	return &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		Hash:         calculateHash(0, time.Now().Unix(), []Transaction{}, "", "GenesisValidator", 0),
		PrevHash:     "",
		Validator:    "GenesisValidator",
		Stake:        0,
	}
}

func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.Chain) > 0 {
		lastBlock := bc.Chain[len(bc.Chain)-1]
		if block.PrevHash != lastBlock.Hash {
			return errors.New("invalid previous hash")
		}
		if block.Index != lastBlock.Index+1 {
			return errors.New("invalid block index")
		}
	}

	// Verify transactions and update account balances
	for _, tx := range block.Transactions {
		if err := bc.verifyAndApplyTransaction(tx); err != nil {
			return err
		}
	}

	bc.Chain = append(bc.Chain, block)
	return nil
}

func (bc *Blockchain) verifyAndApplyTransaction(tx Transaction) error {
	sender, exists := bc.Accounts[tx.From]
	if !exists {
		return errors.New("sender account not found")
	}

	if sender.Balance < tx.Amount+tx.Fee {
		return errors.New("insufficient balance")
	}

	recipient, exists := bc.Accounts[tx.To]
	if !exists {
		recipient = &Account{Address: tx.To, Balance: 0}
		bc.Accounts[tx.To] = recipient
	}

	sender.Balance -= tx.Amount + tx.Fee
	recipient.Balance += tx.Amount

	// Add fee to validator's account (simplified for now)
	validator, exists := bc.Accounts[bc.Chain[len(bc.Chain)-1].Validator]
	if !exists {
		validator = &Account{Address: bc.Chain[len(bc.Chain)-1].Validator, Balance: 0}
		bc.Accounts[bc.Chain[len(bc.Chain)-1].Validator] = validator
	}
	validator.Balance += tx.Fee

	return nil
}

func (bc *Blockchain) GetLastBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.Chain) > 0 {
		return bc.Chain[len(bc.Chain)-1]
	}
	return nil
}

func (bc *Blockchain) GetBalance(address string) (int64, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	account, exists := bc.Accounts[address]
	if !exists {
		return 0, errors.New("account not found")
	}
	return account.Balance, nil
}

func (bc *Blockchain) CreateAccount(name string) (*Account, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	address := generateAddress()
	account := &Account{
		Name:    name,
		Address: address,
		Balance: 0,
		Stake:   0,
	}

	bc.Accounts[address] = account
	return account, nil
}

func (bc *Blockchain) BuyJUL(accountAddress string, qarAmount float64) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	julAmount := int64(qarAmount * float64(ExchangeRate) * NANO)
	if bc.Exchange.Balance < julAmount {
		return errors.New("insufficient JUL in exchange")
	}

	account, exists := bc.Accounts[accountAddress]
	if !exists {
		return errors.New("account not found")
	}

	bc.Exchange.Balance -= julAmount
	account.Balance += julAmount

	return nil
}

func (bc *Blockchain) GetMempool() []Transaction {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.Mempool
}

func (bc *Blockchain) selectValidator() string {
	// Simple stake-weighted random selection
	totalStake := int64(0)
	for _, account := range bc.Accounts {
		totalStake += account.Stake
	}

	if totalStake == 0 {
		// If no stakes, select randomly
		for address := range bc.Accounts {
			return address
		}
	}

	randomStake := time.Now().UnixNano() % totalStake
	cumulativeStake := int64(0)
	for address, account := range bc.Accounts {
		cumulativeStake += account.Stake
		if cumulativeStake > randomStake {
			return address
		}
	}

	// This should never happen, but return the first account if it does
	for address := range bc.Accounts {
		return address
	}
	return ""
}

func calculateHash(index, timestamp int64, transactions []Transaction, prevHash, validator string, stake int64) string {
	record := string(index) + string(timestamp) + string(stake) + prevHash + validator
	for _, tx := range transactions {
		record += tx.From + tx.To + string(tx.Amount) + string(tx.Fee)
	}
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func generateAddress() string {
	// This is a simplified address generation.
	// In a real implementation, you'd use proper cryptographic methods.
	hash := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(hash[:])[:40]
}

type Account struct {
	Name    string
	Address string
	Balance int64
	Stake   int64
}

func (bc *Blockchain) AddStake(address string, amount int64) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	account, exists := bc.Accounts[address]
	if !exists {
		return errors.New("account not found")
	}

	if account.Balance < amount {
		return errors.New("insufficient balance to stake")
	}

	account.Balance -= amount
	account.Stake += amount

	return nil
}

func (bc *Blockchain) RemoveStake(address string, amount int64) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	account, exists := bc.Accounts[address]
	if !exists {
		return errors.New("account not found")
	}

	if account.Stake < amount {
		return errors.New("insufficient stake to remove")
	}

	account.Stake -= amount
	account.Balance += amount

	return nil
}

func (bc *Blockchain) MintInitialSupply(address string, amount int64) error {
	account, exists := bc.Accounts[address]
	if !exists {
		return fmt.Errorf("account %s does not exist", address)
	}
	account.Balance += amount
	return nil
}

func (bc *Blockchain) GetAccounts() []*Account {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	accounts := make([]*Account, 0, len(bc.Accounts))
	for _, account := range bc.Accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

func (bc *Blockchain) GetChain() []*Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.Chain
}
