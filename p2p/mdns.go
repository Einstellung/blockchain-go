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
		if pi.ID > n.h.ID() {
			fmt.Println(utils.ShortID(pi.ID), " id is greater than us, wait for it to connect to us")
		} else {
			err := n.h.Connect(context.Background(), pi)
			if err != nil {
				fmt.Println("error connecting to peer", utils.ShortID(pi.ID), err)
			}
			fmt.Println("connected to peer", utils.ShortID(pi.ID))
		}
}

func SetupDiscovery(h host.Host) error {
	s := mdns.NewMdnsService(h, "blockchain-go", &discoveryNotifee{h: h})
	return s.Start()
}