package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Add these constants at the top of the file
const (
	USDtoJULRate = 35
	BlockReward  = 50 // JUL reward for forging a block
	MinStake     = 100 // Minimum amount of JUL that can be staked (same as MinimumStake in pos.go)
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Balance    float64
	Staked     float64
}

func NewWallet() *Wallet {
	private, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public := &private.PublicKey

	return &Wallet{
		PrivateKey: private,
		PublicKey:  public,
		Balance:    0,
		Staked:     0,
	}
}

func (w *Wallet) GetAddress() string {
	pubKeyBytes := elliptic.Marshal(elliptic.P256(), w.PublicKey.X, w.PublicKey.Y)
	pubKeyHash := sha256.Sum256(pubKeyBytes)
	return hex.EncodeToString(pubKeyHash[:])
}

func (w *Wallet) SignTransaction(tx *Transaction) ([]byte, error) {
	txHash := CalculateTransactionHash(tx)
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

func (w *Wallet) Stake(amount float64) error {
	if amount < MinStake {
		return fmt.Errorf("minimum stake is %d JUL", MinStake)
	}
	if amount > w.Balance {
		return fmt.Errorf("insufficient balance")
	}
	w.Balance -= amount
	w.Staked += amount
	return nil
}

func (w *Wallet) Unstake(amount float64) error {
	if amount > w.Staked {
		return fmt.Errorf("cannot unstake more than staked amount")
	}
	w.Staked -= amount
	w.Balance += amount
	return nil
}

func (w *Wallet) GetStakedAmount() float64 {
	return w.Staked
}

func (w *Wallet) GetTotalBalance() float64 {
	return w.Balance + w.Staked
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

func GetBlockReward() float64 {
	return BlockReward
}

// This function should be defined somewhere in your blockchain package
// func CalculateTransactionHash(tx *Transaction) string {
//     // Implement the hash calculation for a transaction
//     // This is just a placeholder implementation
//     return fmt.Sprintf("%v", tx)
// }