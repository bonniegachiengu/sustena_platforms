package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	//"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func NewWallet() *Wallet {
	private, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public := &private.PublicKey

	return &Wallet{private, public}
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