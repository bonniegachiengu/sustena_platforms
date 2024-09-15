package network

import (
    "context"
    "fmt"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p/core/host"
    "github.com/libp2p/go-libp2p/core/peer"
    "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
    "github.com/multiformats/go-multiaddr"
    "github.com/Sustena_platforms/config"
)

type P2PNetwork struct {
    host host.Host
    ctx  context.Context
}

func NewP2PNetwork(listenAddr string) (*P2PNetwork, error) {
    ctx := context.Background()

    // Check if the multiaddr is empty before using it
    if config.GetNetworkConfig().ListenAddr == "" {
        return nil, fmt.Errorf("Listen address is empty. Please provide a valid multiaddr in the configuration.")
    }

    // Create a libp2p host
    h, err := libp2p.New(
        libp2p.ListenAddrStrings(config.NetworkConfig.ListenAddr),
    )
    if err != nil {
        return nil, fmt.Errorf("Failed to create libp2p host: %w", err)
    }

    p2p := &P2PNetwork{
        host: h,
        ctx:  ctx,
    }

    // Setup local mDNS discovery
    if err := p2p.setupDiscovery(); err != nil {
        return nil, err
    }

    return p2p, nil
}

func (p *P2PNetwork) setupDiscovery() error {
    // Setup local mDNS discovery
    s := mdns.NewMdnsService(p.host, "sustena-network", &discoveryNotifee{p: p})
    return s.Start()
}

func (p *P2PNetwork) Connect(peerAddr string) error {
    maddr, err := multiaddr.NewMultiaddr(peerAddr)
    if err != nil {
        return err
    }

    peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
    if err != nil {
        return err
    }

    if err := p.host.Connect(p.ctx, *peerInfo); err != nil {
        return err
    }

    fmt.Printf("Connected to peer: %s\n", peerInfo.ID.String())
    return nil
}

func (p *P2PNetwork) Broadcast(msg []byte) error {
    // Implement message broadcasting to all connected peers
    // This is a placeholder and should be implemented based on your specific needs
    return nil
}

func (p *P2PNetwork) ReceiveMessage() []byte {
    // Implement message receiving logic
    // This is a placeholder and should be implemented based on your specific needs
    return nil
}

func (p *P2PNetwork) Shutdown() error {
    return p.host.Close()
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
    p *P2PNetwork
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
    fmt.Printf("Discovered new peer %s\n", pi.ID.String())
    err := n.p.host.Connect(n.p.ctx, pi)
    if err != nil {
        fmt.Printf("Error connecting to peer %s: %s\n", pi.ID.String(), err)
    }
}
