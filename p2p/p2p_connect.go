package p2p

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func PeerConnect() {
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
	if err := SetupDiscovery(h); err != nil {
		panic(err)
	}

	// join the blockchain pub/sub network
	bcommu, err := JoinBlockChain(ctx, ps, h.ID())
	if err != nil {
		panic(err)
	}

	// draw the ui
	ui := NewBlockchainUi(bcommu)

	go ui.readDataTerminal()
	go ui.handleEvents()

	defer ui.end()

	// block main thread
	select {}
}
