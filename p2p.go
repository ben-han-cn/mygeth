package mygeth

import (
	"log"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/params"
)

func setP2PConfig(cfg *p2p.Config) {
	setBootstrapNodes(cfg)
	setBootstrapNodesV5(cfg)
	cfg.NoDiscovery = false
	cfg.MaxPeers = 100
	cfg.DiscoveryV5 = false
}

func setBootstrapNodes(cfg *p2p.Config) {
	urls := params.MainnetBootnodes
	cfg.BootstrapNodes = make([]*discover.Node, 0, len(urls))
	for _, url := range urls {
		node, err := discover.ParseNode(url)
		if err != nil {
			log.Printf("Bootstrap URL %s invalid enode %s\n", url, err.Error())
			continue
		}
		cfg.BootstrapNodes = append(cfg.BootstrapNodes, node)
	}
}

func setBootstrapNodesV5(cfg *p2p.Config) {
	urls := params.DiscoveryV5Bootnodes
	cfg.BootstrapNodesV5 = make([]*discv5.Node, 0, len(urls))
	for _, url := range urls {
		node, err := discv5.ParseNode(url)
		if err != nil {
			log.Printf("Bootstrap URL %s invalid enode %s\n", url, err.Error())
			continue
		}
		cfg.BootstrapNodesV5 = append(cfg.BootstrapNodesV5, node)
	}
}
