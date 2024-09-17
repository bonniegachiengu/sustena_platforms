package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sustena_platforms/entropy/blockchain"
	"sustena_platforms/entropy/node"
	"time"
)

type CLI struct {
	Node *node.Node
}

func NewCLI(node *node.Node) *CLI {
	return &CLI{Node: node}
}

func (cli *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nAvailable commands:")
		fmt.Println("1. Create wallet")
		fmt.Println("2. Check balance")
		fmt.Println("3. Send transaction")
		fmt.Println("4. Forge new block")  // Changed from "Mine block"
		fmt.Println("5. Print blockchain")
		fmt.Println("6. Exit")
		fmt.Print("Enter command number: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			cli.createWallet()
		case "2":
			cli.checkBalance()
		case "3":
			cli.sendTransaction()
		case "4":
			cli.forgeBlock()  // Changed from cli.mineBlock()
		case "5":
			cli.printBlockchain()
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command")
		}
	}
}

func (cli *CLI) createWallet() {
	wallet := cli.Node.WalletManager.CreateWallet()
	cli.Node.Blockchain.RegisterWallet(wallet)
	address := wallet.GetAddress()
	fmt.Printf("New wallet created with address: %s\n", address)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter USD amount to purchase JUL: ")
	usdAmountStr, _ := reader.ReadString('\n')
	usdAmountStr = strings.TrimSpace(usdAmountStr)
	usdAmount, err := strconv.ParseFloat(usdAmountStr, 64)
	if err != nil {
		fmt.Println("Invalid amount")
		return
	}

	julPurchased := wallet.PurchaseJUL(usdAmount)
	fmt.Printf("Purchased %.2f JUL\n", julPurchased)

	// Update the blockchain's wallet
	cli.Node.Blockchain.Wallets[address] = wallet
}

func (cli *CLI) checkBalance() {
	fmt.Print("Enter wallet address: ")
	reader := bufio.NewReader(os.Stdin)
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	balance := cli.Node.Blockchain.GetBalance(address)
	fmt.Printf("Balance for %s: %.2f JUL\n", address, balance)
}

func (cli *CLI) sendTransaction() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter sender address: ")
	from, _ := reader.ReadString('\n')
	from = strings.TrimSpace(from)

	fmt.Print("Enter recipient address: ")
	to, _ := reader.ReadString('\n')
	to = strings.TrimSpace(to)

	fmt.Print("Enter amount: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Invalid amount")
		return
	}

	wallet := cli.Node.WalletManager.GetWallet(from)
	if wallet == nil {
		fmt.Println("Sender wallet not found")
		return
	}

	tx, err := blockchain.NewTransaction(wallet, to, amount)
	if err != nil {
		fmt.Printf("Error creating transaction: %v\n", err)
		return
	}

	if cli.Node.Blockchain.AddTransaction(tx) {
		cli.Node.P2P.BroadcastTransaction(cli.Node.ID, tx)
		fmt.Println("Transaction added and broadcasted")
	} else {
		fmt.Println("Failed to add transaction")
	}
}

func (cli *CLI) forgeBlock() {
	validator := cli.Node.Blockchain.PoS.SelectValidator()
	if validator == nil {
		fmt.Println("No validator selected")
		return
	}

	time.Sleep(time.Second) // Add a 1-second delay before creating a new block

	err := cli.Node.Blockchain.AddBlock(validator.Address)
	if err != nil {
		fmt.Printf("Error forging new block: %v\n", err)
	} else {
		newBlock := cli.Node.Blockchain.GetLatestBlock()
		cli.Node.PropagateNewBlock(newBlock)

		// Reward the validator
		validatorWallet := cli.Node.WalletManager.GetWallet(validator.Address)
		if validatorWallet != nil {
			blockReward := blockchain.GetBlockReward() // Use a function to get the block reward
			validatorWallet.AddBalance(blockReward)
			fmt.Printf("Validator %s rewarded with %.2f JUL\n", validator.Address, blockReward)
		}

		fmt.Printf("New block forged and propagated by validator %s\n", validator.Address)
	}
}

func (cli *CLI) printBlockchain() {
	for _, block := range cli.Node.Blockchain.Blocks {
		fmt.Printf("Block %d:\n", block.Index)
		fmt.Printf("  Timestamp: %d\n", block.Timestamp)
		fmt.Printf("  Transactions: %d\n", len(block.Transactions))
		fmt.Printf("  Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("  Hash: %s\n", block.Hash)
		fmt.Printf("  Validator: %s\n", block.Validator)
		fmt.Println()
	}
}