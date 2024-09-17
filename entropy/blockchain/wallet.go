package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	//"math/big"
)

// Add these constants at the top of the file
const (
	USDtoJULRate = 35
	BlockReward  = 50 // JUL reward for forging a block
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Balance    float64 // Add this line
}

func NewWallet() *Wallet {
	private, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public := &private.PublicKey

	return &Wallet{private, public, 0} // Initialize balance to 0
}

func (w *Wallet) GetAddress() string {
	pubKeyBytes := elliptic.Marshal(elliptic.P256(), w.PublicKey.X, w.PublicKey.Y)
	pubKeyHash := sha256.Sum256(pubKeyBytes)
	return hex.EncodeToString(pubKeyHash[:])
}

func (w *Wallet) SignTransaction(tx *Transaction) ([]byte, error) {
	txHash := calculateTransactionHash(tx)
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, []byte(txHash))
	if err != nil {
		return nil, err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	// Pad the signature to ensure it's always 64 bytes
	if len(signature) < 64 {
		padded := make([]byte, 64)
		copy(padded[64-len(signature):], signature)
		signature = padded
	}
	return signature, nil
}

// Add these methods
func (w *Wallet) GetBalance() float64 {
	return w.Balance
}

func (w *Wallet) SetBalance(balance float64) {
	w.Balance = balance
}

func (w *Wallet) AddBalance(amount float64) {
	w.Balance += amount
}

func (w *Wallet) DeductBalance(amount float64) {
	w.Balance -= amount
}

func (w *Wallet) PurchaseJUL(usdAmount float64) float64 {
	julAmount := usdAmount * USDtoJULRate
	w.Balance += julAmount
	return julAmount
}

type WalletManager struct {
	Wallets map[string]*Wallet
}

func NewWalletManager() *WalletManager {
	return &WalletManager{
		Wallets: make(map[string]*Wallet),
	}
}

func (wm *WalletManager) CreateWallet() *Wallet {
	wallet := NewWallet()
	address := wallet.GetAddress()
	wm.Wallets[address] = wallet
	return wallet
}

func (wm *WalletManager) GetWallet(address string) *Wallet {
	return wm.Wallets[address]
}

// Add this function at the end of the file
func GetBlockReward() float64 {
	return BlockReward
}