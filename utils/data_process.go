package utils

import "github.com/libp2p/go-libp2p/core/peer"

func ShortID(peer peer.ID) string {
	pretty := peer.String()
	return pretty[len(pretty)-6:]
}