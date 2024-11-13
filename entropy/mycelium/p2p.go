package mycelium

import (
	"sustena_platforms/entropy/blockchain"
	"sustena_platforms/utils"
	"sync"
	"fmt"
)

type Peer struct {
	ID         string
	Blockchain *blockchain.Blockchain
}

type P2PNetwork struct {
	Peers map[string]*Peer
	mu    sync.Mutex
}

func NewP2PNetwork() *P2PNetwork {
	return &P2PNetwork{
		Peers: make(map[string]*Peer),
	}
}

func (n *P2PNetwork) AddPeer(id string, blockchain *blockchain.Blockchain) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Peers[id] = &Peer{ID: id, Blockchain: blockchain}
}

func (n *P2PNetwork) BroadcastBlock(senderID string, block *blockchain.Block) {
	for id := range n.Peers {
		if id != senderID {
			n.SendBlock(id, block)
		}
	}
}

func (n *P2PNetwork) SendBlock(peerID string, block *blockchain.Block) {
	peer, exists := n.Peers[peerID]
	if !exists {
		err := utils.NewError(fmt.Sprintf("Peer %s not found", peerID))
		utils.LogError(err)
		return
	}

	currentLastBlock := peer.Blockchain.GetLatestBlock()
	if block.Index == currentLastBlock.Index+1 && block.PreviousHash == currentLastBlock.Hash {
		if err := peer.Blockchain.AddBlock(block.Validator); err == nil {
			utils.LogInfo(fmt.Sprintf("Block %d added to the chain of peer %s", block.Index, peerID))
		} else {
			utils.LogError(utils.NewError(fmt.Sprintf("Error adding block %d to peer %s: %v", block.Index, peerID, err)))
		}
	} else if block.Index > currentLastBlock.Index {
		n.RequestChain(peerID)
	}
}

func (n *P2PNetwork) RequestChain(peerID string) {
	peer, exists := n.Peers[peerID]
	if !exists {
		fmt.Printf("Peer %s not found\n", peerID)
		return
	}

	var longestChain []*blockchain.Block
	var longestLength int

	for _, otherPeer := range n.Peers {
		if otherPeer.ID != peerID {
			peerChain := otherPeer.Blockchain.Blocks
			if len(peerChain) > longestLength && peer.Blockchain.IsValidChain(peerChain) {
				longestChain = peerChain
				longestLength = len(peerChain)
			}
		}
	}

	if longestLength > peer.Blockchain.GetChainLength() {
		peer.Blockchain.ReplaceChain(longestChain)
		fmt.Printf("Chain of peer %s replaced with a longer valid chain\n", peerID)
	}
}

func (n *P2PNetwork) BroadcastTransaction(senderID string, tx *blockchain.Transaction) {
	for id := range n.Peers {
		if id != senderID {
			n.SendTransaction(id, tx)
		}
	}
}

func (n *P2PNetwork) SendTransaction(peerID string, tx *blockchain.Transaction) {
	peer, exists := n.Peers[peerID]
	if !exists {
		fmt.Printf("Peer %s not found\n", peerID)
		return
	}

	if peer.Blockchain.AddTransaction(tx) {
		fmt.Printf("Transaction %s added to the mempool of peer %s\n", tx.ID, peerID)
	} else {
		fmt.Printf("Failed to add transaction %s to the mempool of peer %s\n", tx.ID, peerID)
	}
}

func (n *P2PNetwork) SyncBlockchain(peerID string) {
	peer, exists := n.Peers[peerID]
	if !exists {
		utils.LogError(utils.NewError(fmt.Sprintf("Peer %s not found", peerID)))
		return
	}

	longestChain := peer.Blockchain.Blocks
	longestLength := len(longestChain)

	for _, otherPeer := range n.Peers {
		if otherPeer.ID != peerID {
			peerChain := otherPeer.Blockchain.Blocks
			if len(peerChain) > longestLength && peer.Blockchain.IsValidChain(peerChain) {
				longestChain = peerChain
				longestLength = len(peerChain)
			}
		}
	}

	if longestLength > peer.Blockchain.GetChainLength() {
		err := peer.Blockchain.ReplaceChain(longestChain)
		if err != nil {
			utils.LogError(utils.NewError(fmt.Sprintf("Failed to replace chain for peer %s: %v", peerID, err)))
		} else {
			utils.LogInfo(fmt.Sprintf("Blockchain of peer %s synced with the longest valid chain (length: %d)", peerID, longestLength))
		}
	} else {
		utils.LogInfo(fmt.Sprintf("Blockchain of peer %s is already up to date (length: %d)", peerID, peer.Blockchain.GetChainLength()))
	}
}
