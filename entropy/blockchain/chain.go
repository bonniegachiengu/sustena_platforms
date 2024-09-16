package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

type Account struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
	Balance    int64 // Balance in Nanojoules
	Nonce      uint64
}

type Blockchain struct {
	Chain    []*Block
	Accounts map[string]*Account
	mu       sync.RWMutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain:    []*Block{GenesisBlock()},
		Accounts: make(map[string]*Account),
	}
}

func GenesisBlock() *Block {
	return NewBlock(0, []Transaction{}, "", "GenesisValidator", 0)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	if len(bc.Chain) > 0 {
		lastBlock := bc.Chain[len(bc.Chain)-1]
		if block.PrevHash != lastBlock.Hash {
			return errors.New("invalid previous hash")
		}
		if block.Index != lastBlock.Index+1 {
			return errors.New("invalid block index")
		}
	}

	bc.Chain = append(bc.Chain, block)
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
	account, err := bc.GetAccount(address)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

// Add this method to your Blockchain struct
func (bc *Blockchain) DeployContract(bytecode []byte) (string, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Generate a simple contract address (this is a simplified example)
	hash := sha256.Sum256(bytecode)
	address := hex.EncodeToString(hash[:])

	// In a real implementation, you'd store the contract somewhere
	// For now, we'll just return the address
	return address, nil
}

func (bc *Blockchain) CreateAccount() (*Account, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey
	address := generateAddress(&publicKey)

	account := &Account{
		Address:    address,
		PrivateKey: privateKey,
		Balance:    0,
		Nonce:      0,
	}

	bc.mu.Lock()
	bc.Accounts[address] = account
	bc.mu.Unlock()

	return account, nil
}

func (bc *Blockchain) GetAccount(address string) (*Account, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	account, exists := bc.Accounts[address]
	if !exists {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (bc *Blockchain) Transfer(from, to string, amount int64) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	fromAccount, exists := bc.Accounts[from]
	if !exists {
		return errors.New("sender account not found")
	}

	if fromAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	toAccount, exists := bc.Accounts[to]
	if !exists {
		return errors.New("recipient account not found")
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount

	tx := Transaction{
		From:   from,
		To:     to,
		Amount: amount,
		Nonce:  fromAccount.Nonce,
	}

	fromAccount.Nonce++

	var lastBlock *Block
	if len(bc.Chain) > 0 {
		lastBlock = bc.Chain[len(bc.Chain)-1]
	} else {
		return errors.New("blockchain is empty")
	}

	newBlock := NewBlock(lastBlock.Index+1, []Transaction{tx}, lastBlock.Hash, "Validator", 0) // For simplicity, we're not implementing PoS yet
	return bc.AddBlock(newBlock)
}

func generateAddress(publicKey *ecdsa.PublicKey) string {
	pubBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	hash := sha256.Sum256(pubBytes)
	return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) MintInitialSupply(address string, amount int64) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	account, exists := bc.Accounts[address]
	if !exists {
		return errors.New("account not found")
	}

	account.Balance += amount

	tx := Transaction{
		From:   "MintingAddress",
		To:     address,
		Amount: amount,
		Nonce:  0,
	}

	var lastBlock *Block
	if len(bc.Chain) > 0 {
		lastBlock = bc.Chain[len(bc.Chain)-1]
	} else {
		return errors.New("blockchain is empty")
	}

	newBlock := NewBlock(lastBlock.Index+1, []Transaction{tx}, lastBlock.Hash, "Validator", 0)
	return bc.AddBlock(newBlock)
}
