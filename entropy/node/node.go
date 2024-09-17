package node

import (
	"fmt"
	"net/http"
	"sustena_platforms/entropy/blockchain"
	"sustena_platforms/entropy/mycelium"
	"sustena_platforms/utils"
)

type Node struct {
	ID            string
	Port          int
	Blockchain    *blockchain.Blockchain
	P2P           *mycelium.P2PNetwork
	WalletManager *blockchain.WalletManager
	mux           *http.ServeMux
}

func NewNode(id string, port int, stateFile string) *Node {
	bc := blockchain.NewBlockchain(stateFile)
	p2p := mycelium.NewP2PNetwork()
	return &Node{
		ID:         id,
		Port:       port,
		Blockchain: bc,
		P2P:        p2p,
		mux:        http.NewServeMux(),
	}
}

func (n *Node) SetWalletManager(wm *blockchain.WalletManager) {
	n.WalletManager = wm
}

func (n *Node) SetP2PNetwork(network *mycelium.P2PNetwork) {
	n.P2P = network
}

func (n *Node) Start() {
	n.setupHTTPHandlers()
	go n.startHTTPServer()
	fmt.Printf("Node %s started on port %d\n", n.ID, n.Port)
}

func (n *Node) setupHTTPHandlers() {
	n.mux.HandleFunc("/status", n.handleStatus)
	n.mux.HandleFunc("/sync", n.handleSync)
	n.mux.HandleFunc("/transaction", n.handleTransaction)
}

func (n *Node) startHTTPServer() {
	http.ListenAndServe(fmt.Sprintf(":%d", n.Port), n.mux)
}

func (n *Node) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := fmt.Sprintf("Node %s\nBlockchain length: %d\nLatest block hash: %s\n",
		n.ID, n.Blockchain.GetChainLength(), n.Blockchain.GetLatestBlock().Hash)
	w.Write([]byte(status))
}

func (n *Node) handleSync(w http.ResponseWriter, r *http.Request) {
	n.P2P.SyncBlockchain(n.ID)
	w.Write([]byte("Blockchain sync initiated"))
}

func (n *Node) handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := utils.NewError("Only POST method is allowed")
		utils.LogError(err)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	wallet := n.WalletManager.CreateWallet()
	tx, err := blockchain.NewTransaction(wallet, "recipient_address", 10)
	if err != nil {
		err = utils.NewError(fmt.Sprintf("Error creating transaction: %v", err))
		utils.LogError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if n.Blockchain.AddTransaction(tx) {
		n.P2P.BroadcastTransaction(n.ID, tx)
		utils.LogInfo(fmt.Sprintf("Transaction %s created and broadcasted", tx.ID))
		w.Write([]byte(fmt.Sprintf("Transaction %s created and broadcasted", tx.ID)))
	} else {
		err := utils.NewError("Failed to add transaction to mempool")
		utils.LogError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (n *Node) PropagateNewBlock(block *blockchain.Block) {
	n.P2P.BroadcastBlock(n.ID, block)
}