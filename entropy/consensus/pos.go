package consensus

import (
    "math/rand"
    "time"
    "github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
)

type ProofOfStake struct {
    Validators map[string]int64
}

func NewProofOfStake() *ProofOfStake {
    return &ProofOfStake{
        Validators: make(map[string]int64),
    }
}

func (pos *ProofOfStake) AddValidator(address string, stake int64) {
    pos.Validators[address] = stake
}

func (pos *ProofOfStake) RemoveValidator(address string) {
    delete(pos.Validators, address)
}

func (pos *ProofOfStake) SelectValidator() string {
    totalStake := int64(0)
    for _, stake := range pos.Validators {
        totalStake += stake
    }

    source := rand.NewSource(time.Now().UnixNano())
    r := rand.New(source)
    randomNumber := r.Int63n(totalStake)
    
    cumulativeStake := int64(0)
    for validator, stake := range pos.Validators {
        cumulativeStake += stake
        if randomNumber < cumulativeStake {
            return validator
        }
    }

    // This should never happen, but return the first validator if it does
    for validator := range pos.Validators {
        return validator
    }
    return ""
}

func (pos *ProofOfStake) Validate(block *blockchain.Block) bool {
    if _, exists := pos.Validators[block.Validator]; !exists {
        return false
    }

    if pos.Validators[block.Validator] < block.Stake {
        return false
    }

    // Additional validation logic can be added here

    return true
}
