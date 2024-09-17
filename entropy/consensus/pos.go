package consensus

import (
	"fmt"
	"math/rand"
	"time"
	"sustena_platforms/utils"
)

const (
	MinimumStake     = 100.0 // Minimum stake in JUL
	SlashingPercent  = 0.05  // 5% slashing for violations
	CommunityFundAddress = "community_fund_address" // Replace with actual address
)

// This is a placeholder for the Proof of Stake consensus mechanism
// We'll implement a simple version for now and expand later

type Validator struct {
	Address string
	Stake   float64
	TotalBalance float64
}

type PoS struct {
	Validators     []Validator
	TotalStake     float64
	LastValidation time.Time
	CommunityFund  float64
}

func NewPoS() *PoS {
	return &PoS{
		Validators:     []Validator{},
		TotalStake:     0,
		LastValidation: time.Now(),
		CommunityFund:  0,
	}
}

func (pos *PoS) AddValidator(address string, stake float64, totalBalance float64) {
	if stake >= MinimumStake {
		pos.Validators = append(pos.Validators, Validator{Address: address, Stake: stake, TotalBalance: totalBalance})
		pos.TotalStake += stake
		utils.LogInfo(fmt.Sprintf("Added validator %s with stake %.2f and total balance %.2f", address, stake, totalBalance))
	} else {
		utils.LogInfo(fmt.Sprintf("Validator %s not added: stake %.2f is below minimum stake %.2f", address, stake, MinimumStake))
	}
}

func (pos *PoS) UpdateStake(address string, newStake float64, newTotalBalance float64) {
	for i, validator := range pos.Validators {
		if validator.Address == address {
			oldStake := validator.Stake
			pos.TotalStake = pos.TotalStake - validator.Stake + newStake
			pos.Validators[i].Stake = newStake
			pos.Validators[i].TotalBalance = newTotalBalance
			utils.LogInfo(fmt.Sprintf("Updated stake for validator %s: %.2f -> %.2f, new total balance: %.2f", address, oldStake, newStake, newTotalBalance))
			if newStake < MinimumStake {
				pos.RemoveValidator(address)
			}
			return
		}
	}
	utils.LogInfo(fmt.Sprintf("Validator %s not found for stake update", address))
}

func (pos *PoS) RemoveValidator(address string) {
	for i, validator := range pos.Validators {
		if validator.Address == address {
			pos.TotalStake -= validator.Stake
			pos.Validators = append(pos.Validators[:i], pos.Validators[i+1:]...)
			utils.LogInfo(fmt.Sprintf("Removed validator %s with stake %.2f", address, validator.Stake))
			return
		}
	}
	utils.LogInfo(fmt.Sprintf("Validator %s not found for removal", address))
}

func (pos *PoS) SelectValidator() *Validator {
	utils.LogInfo(fmt.Sprintf("Selecting validator. Total validators: %d", len(pos.Validators)))
	if len(pos.Validators) == 0 {
		fmt.Println("No validators available")
		return nil
	}

	// Remove the time constraint for now
	// if time.Since(pos.LastValidation) < 5*time.Second {
	//     return nil
	// }

	// Calculate selection weights
	weights := make([]float64, len(pos.Validators))
	totalWeight := 0.0

	for i, validator := range pos.Validators {
		// Weight is based on stake amount, time since last selection, and staked percentage
		timeFactor := time.Since(pos.LastValidation).Hours() / 24 // Convert to days
		stakePercentage := validator.Stake / validator.TotalBalance
		weight := validator.Stake * (1 + timeFactor) * (1 + stakePercentage)
		weights[i] = weight
		totalWeight += weight
		fmt.Printf("Validator %s: Stake: %.2f, Weight: %.2f\n", validator.Address, validator.Stake, weight)
	}

	fmt.Printf("Total weight: %.2f\n", totalWeight)

	// Select a validator based on weights
	randomValue := rand.Float64() * totalWeight
	fmt.Printf("Random value: %.2f\n", randomValue)
	for i, weight := range weights {
		randomValue -= weight
		if randomValue <= 0 {
			pos.LastValidation = time.Now()
			fmt.Printf("Selected validator: %s\n", pos.Validators[i].Address)
			return &pos.Validators[i]
		}
	}

	// This should never happen, but just in case
	pos.LastValidation = time.Now()
	fmt.Printf("Fallback: Selected last validator: %s\n", pos.Validators[len(pos.Validators)-1].Address)
	return &pos.Validators[len(pos.Validators)-1]
}

func (pos *PoS) SlashValidator(address string) float64 {
	for i, validator := range pos.Validators {
		if validator.Address == address {
			slashedAmount := validator.Stake * SlashingPercent
			pos.Validators[i].Stake -= slashedAmount
			pos.TotalStake -= slashedAmount
			pos.CommunityFund += slashedAmount
			utils.LogInfo(fmt.Sprintf("Slashed validator %s: %.2f JUL, new stake: %.2f JUL", address, slashedAmount, pos.Validators[i].Stake))
			if pos.Validators[i].Stake < MinimumStake {
				pos.RemoveValidator(address)
			}
			return slashedAmount
		}
	}
	utils.LogInfo(fmt.Sprintf("Validator %s not found for slashing", address))
	return 0
}

func (pos *PoS) GetCommunityFund() float64 {
	return pos.CommunityFund
}
