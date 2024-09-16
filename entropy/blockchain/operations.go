package blockchain

import (
    "fmt"
)

type BlockchainOperations struct {
    BC *Blockchain
}

func NewBlockchainOperations(bc *Blockchain) *BlockchainOperations {
    return &BlockchainOperations{BC: bc}
}

func (ops *BlockchainOperations) CreateNewAccount() (string, error) {
    account, err := ops.BC.CreateAccount()
    if err != nil {
        return "", err
    }
    return account.Address, nil
}

func (ops *BlockchainOperations) GetAccountBalance(address string) (float64, error) {
    balance, err := ops.BC.GetBalance(address)
    if err != nil {
        return 0, err
    }
    return float64(balance) / JouleToNanojoule, nil
}

func (ops *BlockchainOperations) TransferJoules(from, to string, amount float64) error {
    amountNanojoules := int64(amount * JouleToNanojoule)
    return ops.BC.Transfer(from, to, amountNanojoules)
}

func (ops *BlockchainOperations) MintInitialSupply(address string, amount float64) error {
    amountNanojoules := int64(amount * JouleToNanojoule)
    return ops.BC.MintInitialSupply(address, amountNanojoules)
}

func (ops *BlockchainOperations) PrintBlockchainState() {
    fmt.Println("Blockchain State:")
    for _, block := range ops.BC.Chain {
        fmt.Printf("Block %d:\n", block.Index)
        fmt.Printf("  Timestamp: %d\n", block.Timestamp)
        fmt.Printf("  Transactions: %d\n", len(block.Transactions))
        fmt.Printf("  PrevHash: %s\n", block.PrevHash)
        fmt.Printf("  Hash: %s\n", block.Hash)
        fmt.Printf("  Validator: %s\n", block.Validator)
        fmt.Printf("  Stake: %d\n", block.Stake)
        fmt.Println()
    }

    fmt.Println("Account Balances:")
    for address, account := range ops.BC.Accounts {
        fmt.Printf("  %s: %.6f Joules\n", address, float64(account.Balance)/JouleToNanojoule)
    }
}
