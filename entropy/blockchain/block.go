package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        int64
	Timestamp    int64
	Transactions []*Transaction
	PreviousHash string
	Hash         string
	Validator    string
	// Remove the Reward field
}

func CreateBlock(index int64, transactions []*Transaction, previousHash string, validator string) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PreviousHash: previousHash,
		Validator:    validator,
		// Remove the Reward calculation
	}
	block.Hash = CalculateHash(block)
	return block
}

func CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%d%v%s%s", block.Index, block.Timestamp, block.Transactions, block.PreviousHash, block.Validator)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (b *Block) IsValid(previousBlock *Block) bool {
	if previousBlock != nil {
		if b.Index != previousBlock.Index+1 {
			return false
		}
		if b.PreviousHash != previousBlock.Hash {
			return false
		}
		if b.Timestamp <= previousBlock.Timestamp {
			return false
		}
	}
	if CalculateHash(*b) != b.Hash {
		return false
	}
	return true
}

func (b *Block) HasValidTransactions() bool {
	transactionHashes := []string{}
	for _, tx := range b.Transactions {
		publicKey := &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     tx.PublicKey.X,
			Y:     tx.PublicKey.Y,
		}
		if !VerifyTransaction(tx, publicKey) {
			return false
		}
		if contains(transactionHashes, tx.ID) {
			return false // Double-spending attempt
		}
		transactionHashes = append(transactionHashes, tx.ID)
	}
	return true
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
