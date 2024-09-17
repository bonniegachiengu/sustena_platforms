package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

const (
	JUL            = 1
	MICRO          = JUL / 1000000
	NANO           = JUL / 1000000000
	PICO           = JUL / 1000000000000
	BlockTime      = 30 * time.Second
	ExchangeRate   = 35 // 35 JUL to 1 QAR
	InitialSupply  = 1000000 * JUL
)

type Block struct {
	Index        int64
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Validator    string
	Stake        int64
}

type Transaction struct {
	From   string
	To     string
	Amount int64 // Amount in Nanojoules
	Nonce  uint64
	Fee    int64
}

func NewBlock(index int64, transactions []Transaction, prevHash string, validator string, stake int64) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Validator:    validator,
		Stake:        stake,
	}
	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() string {
	record := strconv.FormatInt(b.Index, 10) + strconv.FormatInt(b.Timestamp, 10) + b.PrevHash + b.Validator + strconv.FormatInt(b.Stake, 10)
	for _, tx := range b.Transactions {
		record += tx.From + tx.To + strconv.FormatInt(tx.Amount, 10)
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
