package p2p

import (
	"blockchain-go/utils"
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// the number of incoming messages to buffer for each topic.
const CommunicateBufferSize = 128

const blockchainTopic = "beacon-chain"

type BlockMessage struct {
	Message  string
	SenderID peer.ID
}

type BlockchainCommunication struct {
	ctx      context.Context
	Messages chan BlockMessage
	sub      *pubsub.Subscription
	ps       *pubsub.PubSub
	topic    *pubsub.Topic

	// for the sake of discard self broadcast message
	self peer.ID
}

func JoinBlockChain(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID) (*BlockchainCommunication, error) {
	topic, err := ps.Join(blockchainTopic)
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	bcommu := &BlockchainCommunication{
		ctx:      ctx,
		Messages: make(chan BlockMessage, CommunicateBufferSize),
		sub:      sub,
		ps:       ps,
		topic:    topic,

		self:     selfID,
	}

	// start reading messages from subscription in a loop
	go bcommu.readLoop()
	return bcommu, nil
}

func (bcommu *BlockchainCommunication) Publish(message string) error {
	m := BlockMessage{
		Message:  message,
		SenderID: bcommu.self,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return bcommu.topic.Publish(bcommu.ctx, msgBytes)
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (bcommu *BlockchainCommunication) readLoop() {
	for {
		msg, err := bcommu.sub.Next(bcommu.ctx)
		if err != nil {
			close(bcommu.Messages)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == bcommu.self {
			continue
		}

		bm := new(BlockMessage)
		err = json.Unmarshal(msg.Data, bm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		bcommu.Messages <- *bm
	}
}

func (bcommu *BlockchainCommunication) ListPeers() []peer.ID{
	return bcommu.ps.ListPeers(blockchainTopic)
}

func (bcommu *BlockchainCommunication) SelfID() string {
	return utils.ShortID(bcommu.self)
}