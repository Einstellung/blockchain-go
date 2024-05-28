package cli

import (
	p2p "blockchain-go/p2p"
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (cli *CLI) peerConnect() {
	ctx := context.Background()

	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(
		// Use the keypair we generated
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		panic(err)
	}

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		log.Fatalf("Failed to create PubSub service: %v", err)
	}

	// setup local mDNS discovery
	if err := p2p.SetupDiscovery(h); err != nil {
		panic(err)
	}

	// join the blockchain pub/sub network
	bcommu, err := p2p.JoinBlockChain(ctx, ps, h.ID())
	if err != nil {
		panic(err)
	}

	// draw the ui
	ui := NewBlockchainUi(bcommu)
	if err = ui.Run(); err != nil {
		fmt.Println("error running text UI:", err)
	}
}
