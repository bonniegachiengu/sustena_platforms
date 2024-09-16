package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
	"github.com/bonniegachiengu/sustena_platforms/entropy/consensus"
	"github.com/bonniegachiengu/sustena_platforms/entropy/network"
	"github.com/bonniegachiengu/sustena_platforms/symmetry/interpreter"
	"github.com/bonniegachiengu/sustena_platforms/symmetry/vm"
	"github.com/bonniegachiengu/sustena_platforms/embroidery/compiler"
	"github.com/bonniegachiengu/sustena_platforms/embroidery/parser"
	"github.com/bonniegachiengu/sustena_platforms/config"
)

func main() {
	fmt.Println("Welcome to Sustena Platforms CLI")

	// Initialize components
	cfg, _ := config.LoadConfig()
	bc := blockchain.NewBlockchain()
	_ = consensus.NewProofOfStake()
	_, _ = network.NewP2PNetwork(cfg.NetworkConfig.ListenAddr)
	interpreter := interpreter.NewInterpreter()
	_ = vm.NewVM()
	_ = compiler.NewCompiler()
	_ = parser.NewParser()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {
		case "help":
			printHelp()
		case "balance":
			if len(parts) < 2 {
				fmt.Println("Usage: balance <address>")
				continue
			}
			balance, err := bc.GetBalance(parts[1])
			if err != nil {
				fmt.Printf("Error getting balance: %v\n", err)
				continue
			}
			fmt.Printf("Balance of %s: %d\n", parts[1], balance)
		case "send":
			if len(parts) < 4 {
				fmt.Println("Usage: send <from> <to> <amount>")
				continue
			}
			// Implement send transaction logic
			fmt.Println("Transaction sent (not implemented)")
		case "deploy":
			if len(parts) < 2 {
				fmt.Println("Usage: deploy <contract_code>")
				continue
			}
			// Implement contract deployment logic
			fmt.Println("Contract deployed (not implemented)")
		case "run":
			if len(parts) < 2 {
				fmt.Println("Usage: run <symmetry_code>")
				continue
			}
			code := strings.Join(parts[1:], " ")
			err := interpreter.Interpret(code)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Code executed successfully")
			}
		case "peers":
			// Implement peer listing logic
			fmt.Println("Connected peers (not implemented)")
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}

	fmt.Println("Goodbye!")
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                - Show this help message")
	fmt.Println("  balance <address>   - Get balance of an address")
	fmt.Println("  send <from> <to> <amount> - Send transaction")
	fmt.Println("  deploy <contract_code> - Deploy a smart contract")
	fmt.Println("  run <symmetry_code> - Run Symmetry code")
	fmt.Println("  peers               - List connected peers")
	fmt.Println("  exit                - Exit the CLI")
}
