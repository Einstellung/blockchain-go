package p2p

import (
	"blockchain-go/utils"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Println("discovered new peer", utils.ShortID(pi.ID))
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		utils.PrintErr("error connecting to peer %s", utils.ShortID(pi.ID))
		// utils.PrintErr("error: %s", err)
	}
}

func SetupDiscovery(h host.Host) error {
	s := mdns.NewMdnsService(h, "blockchain-go", &discoveryNotifee{h: h})
	return s.Start()
}