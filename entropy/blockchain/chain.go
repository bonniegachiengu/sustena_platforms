package blockchain

import (
	"errors"
	"sync"
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
