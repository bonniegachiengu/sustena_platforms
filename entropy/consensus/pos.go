package consensus

import (
	"math/rand"
	"fmt"
	// "time"
)

// This is a placeholder for the Proof of Stake consensus mechanism
// We'll implement a simple version for now and expand later

type Validator struct {
	Address string
	Stake   uint64
}

type PoS struct {
	Validators []Validator
	TotalStake uint64
}

func NewPoS() *PoS {
	return &PoS{
		Validators: []Validator{},
		TotalStake: 0,
	}
}

func (pos *PoS) AddValidator(address string, stake uint64) {
	pos.Validators = append(pos.Validators, Validator{Address: address, Stake: stake})
	pos.TotalStake += stake
}

func (pos *PoS) SelectValidator() *Validator {
	if len(pos.Validators) == 0 {
		fmt.Println("No validators available")
		return nil
	}

	totalStake := pos.TotalStake
	if totalStake == 0 {
		fmt.Println("Total stake is 0")
		return nil
	}

	randomStake := uint64(rand.Int63n(int64(totalStake)))
	fmt.Printf("Random stake: %d, Total stake: %d\n", randomStake, totalStake)

	var cumulativeStake uint64
	for _, validator := range pos.Validators {
		cumulativeStake += validator.Stake
		if randomStake < cumulativeStake {
			fmt.Printf("Selected validator: %s with stake %d\n", validator.Address, validator.Stake)
			return &validator
		}
	}

	fmt.Println("No validator selected, returning last validator")
	return &pos.Validators[len(pos.Validators)-1]
}

func (pos *PoS) UpdateStake(address string, newStake uint64) {
	for i, validator := range pos.Validators {
		if validator.Address == address {
			pos.TotalStake = pos.TotalStake - validator.Stake + newStake
			pos.Validators[i].Stake = newStake
			break
		}
	}
}
