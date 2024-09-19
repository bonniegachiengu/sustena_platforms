package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

const TransactionFeePercentage = 0.001 // 0.1% fee

type Transaction struct {
	ID        string
	From      string
	To        string
	Amount    float64
	Fee       float64 // Add this field
	Timestamp int64
	Signature []byte
	PublicKey *PublicKeyJSON
}

type PublicKeyJSON struct {
	X, Y *big.Int
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	type TransactionAlias Transaction
	return json.Marshal(&struct {
		*TransactionAlias
		PublicKey *PublicKeyJSON `json:"PublicKey"`
	}{
		TransactionAlias: (*TransactionAlias)(tx),
		PublicKey:        &PublicKeyJSON{tx.PublicKey.X, tx.PublicKey.Y},
	})
}

func (tx *Transaction) UnmarshalJSON(data []byte) error {
	type TransactionAlias Transaction
	aux := &struct {
		*TransactionAlias
		PublicKey *PublicKeyJSON `json:"PublicKey"`
	}{
		TransactionAlias: (*TransactionAlias)(tx),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	tx.PublicKey = aux.PublicKey
	return nil
}

func NewTransaction(wallet *Wallet, to string, amount float64) (*Transaction, error) {
	if wallet == nil {
		return nil, fmt.Errorf("wallet is nil")
	}
	fee := amount * TransactionFeePercentage
	tx := &Transaction{
		From:      wallet.GetAddress(),
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
		PublicKey: &PublicKeyJSON{X: wallet.PublicKey.X, Y: wallet.PublicKey.Y},
	}
	tx.ID = CalculateTransactionHash(tx)
	signature, err := wallet.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	tx.Signature = signature
	return tx, nil
}

// Ensure this function is exported (starts with a capital letter)
func CalculateTransactionHash(tx *Transaction) string {
	record := fmt.Sprintf("%s%s%f%f%d", tx.From, tx.To, tx.Amount, tx.Fee, tx.Timestamp)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func VerifyTransaction(tx *Transaction, pubKey *ecdsa.PublicKey) bool {
	if tx == nil || pubKey == nil || len(tx.Signature) < 64 {
		return false
	}
	txHash := CalculateTransactionHash(tx)
	r := new(big.Int).SetBytes(tx.Signature[:32])
	s := new(big.Int).SetBytes(tx.Signature[32:])
	return ecdsa.Verify(pubKey, []byte(txHash), r, s)
}