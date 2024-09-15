package blockchain

import (
	"errors"
	"sync"
	"crypto/sha256"
	"encoding/hex"
)

type Blockchain struct {
	Chain  []*Block
	mu     sync.Mutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []*Block{GenesisBlock()},
	}
}

func GenesisBlock() *Block {
	return NewBlock(0, []Transaction{}, "", "GenesisValidator", 0)
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

	bc.Chain = append(bc.Chain, block)
	return nil
}

func (bc *Blockchain) GetLastBlock() *Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.Chain) > 0 {
		return bc.Chain[len(bc.Chain)-1]
	}
	return nil
}

func (bc *Blockchain) GetBalance(address string) int64 {
	balance := int64(0)
	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.From == address {
				balance -= tx.Amount
			}
			if tx.To == address {
				balance += tx.Amount
			}
		}
	}
	return balance
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
