package main

import (
	"fmt"
	"log"

	"github.com/bonniegachiengu/sustena_platforms/entropy/blockchain"
	"github.com/bonniegachiengu/sustena_platforms/entropy/consensus"
	"github.com/bonniegachiengu/sustena_platforms/entropy/network"
	"github.com/bonniegachiengu/sustena_platforms/symmetry/interpreter"
	"github.com/bonniegachiengu/sustena_platforms/symmetry/vm"
	"github.com/bonniegachiengu/sustena_platforms/embroidery/compiler"
	"github.com/bonniegachiengu/sustena_platforms/embroidery/parser"
	"github.com/bonniegachiengu/sustena_platforms/api"
	"github.com/bonniegachiengu/sustena_platforms/utils"
	"github.com/bonniegachiengu/sustena_platforms/config"
)

func main() {
	fmt.Println("Starting Sustena Platform")
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize blockchain
	bc := blockchain.NewBlockchain()

	// Initialize Proof of Stake consensus
	pos := consensus.NewProofOfStake()

	// Initialize P2P network
	p2p, err := network.NewP2PNetwork(cfg.NetworkConfig.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to initialize P2P network: %v", err)
	}
	defer p2p.Shutdown()

	// Connect to bootstrap peers
	for _, peerAddr := range cfg.NetworkConfig.BootstrapPeers {
		if err := p2p.Connect(peerAddr); err != nil {
			log.Printf("Failed to connect to peer %s: %v", peerAddr, err)
		}
	}

	// Initialize Symmetry components
	interpreter := interpreter.NewInterpreter()
	vm := vm.NewVM()

	//Run Symmetry
	runSymmetryScripts(interpreter, vm)

	// Initialize Embroidery components
	compiler := compiler.NewCompiler()
	parser := parser.NewParser()

	// Initialize API server
	apiServer := api.NewServer(api.APIConfig{
		// Map the fields from cfg.APIConfig to api.APIConfig
		// For example:
		// Port: cfg.APIConfig.Port,
		// Host: cfg.APIConfig.Host,
		// ... other fields ...
	})
	go apiServer.Start()
	defer apiServer.Shutdown()

	// Start the main application loop
	for {
		// Handle incoming P2P messages
		msg := p2p.ReceiveMessage()
		if msg != nil {
			handleMessage(msg, bc, pos, interpreter, vm, compiler, parser)
		}

		// Process pending transactions
		processPendingTransactions(bc, pos)

		// Run Symmetry scripts
		runSymmetryScripts(interpreter, vm)

		// Compile and deploy Embroidery contracts
		compileAndDeployContracts(compiler, parser, bc)

		// Check for exit condition
		if shouldExit() {
			break
		}
	}

	// Cleanup and shutdown
	apiServer.Shutdown()
	p2p.Shutdown()
	utils.Cleanup()

	fmt.Println("Sustena Platform shutdown complete")
}

func handleMessage(msg []byte, bc *blockchain.Blockchain, pos *consensus.ProofOfStake, 
	interpreter *interpreter.Interpreter, vm *vm.VM, compiler *compiler.Compiler, parser *parser.Parser) {
	// Handle different types of messages
	// This is a placeholder and should be implemented based on your message types
}

func processPendingTransactions(bc *blockchain.Blockchain, pos *consensus.ProofOfStake) {
	// Process any pending transactions and create new blocks
	// This is a placeholder and should be implemented based on your transaction processing logic
}

func runSymmetryScripts(interpreter *interpreter.Interpreter, vm *vm.VM) {
	// Example usage
	code := `
		x = 10
		y = 20
		z = x + y
	`
	err := interpreter.Interpret(code)
	if err != nil {
		log.Printf("Error interpreting Symmetry script: %v", err)
		return
	}

	
	// Print the results
	memory := vm.GetMemory()
	for key, value := range memory {
		fmt.Printf("%s = %v\n", key, value)
	}
}

func compileAndDeployContracts(compiler *compiler.Compiler, parser *parser.Parser, bc *blockchain.Blockchain) {
	// Example Embroidery contract code
	contractCode := `
		contract SimpleStorage {
			uint256 storedData;

			function set(uint256 x) public {
				storedData = x;
			}

			function get() public view returns (uint256) {
				return storedData;
			}
		}
	`

	// Parse the contract code
	ast, err := parser.Parse(contractCode)
	if err != nil {
		log.Printf("Error parsing Embroidery contract: %v", err)
		return
	}

	// Compile the AST to bytecode
	bytecode, err := compiler.Compile(ast)
	if err != nil {
		log.Printf("Error compiling Embroidery contract: %v", err)
		return
	}

	// Deploy the contract to the blockchain
	contractAddress, err := bc.DeployContract(bytecode)
	if err != nil {
		log.Printf("Error deploying Embroidery contract: %v", err)
		return
	}

	log.Printf("Embroidery contract deployed at address: %s", contractAddress)
}

func shouldExit() bool {
	// Check if the application should exit
	// This is a placeholder and should be implemented based on your exit conditions
	return false
}
