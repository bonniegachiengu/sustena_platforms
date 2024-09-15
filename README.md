# Sustena Platforms

Sustena Platforms is a comprehensive blockchain-based system that integrates multiple components to create a robust and flexible decentralized platform.

## Components

1. **Entropy**: The blockchain core
   - Implements a Proof of Stake consensus mechanism
   - Manages the blockchain and transactions

2. **Symmetry**: A custom virtual machine and interpreter
   - Executes custom scripts and smart contracts

3. **Embroidery**: A smart contract language and compiler
   - Allows for writing and deploying custom smart contracts

4. **P2P Network**: Decentralized peer-to-peer communication
   - Utilizes libp2p for robust networking capabilities

## Getting Started

### Prerequisites

- Go 1.22.0 or later

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/bonniegachiengu/sustena_platforms.git
   ```

2. Navigate to the project directory:
   ```bash
   cd sustena_platforms
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Configuration

The project uses a YAML configuration file located at `config/config.yaml`. Modify this file to adjust network and API settings.

### Running the Project

To start the Sustena Platform:

```
go run main.go
```

## Features

- Blockchain implementation with Proof of Stake consensus
- Custom virtual machine for executing smart contracts
- P2P networking for decentralized communication
- API server for interacting with the platform
- Configuration management using Viper

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

