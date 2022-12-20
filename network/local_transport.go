package network

import (
	"bytes"
	"fmt"
	"sync"
)

// transport responsable for maintainng and connecting to peers
type LocalTransport struct {
	address        NetAddress
	consumeChannel chan RPC
	lock           sync.RWMutex
	peers          map[NetAddress]*LocalTransport
}

func NewLocalTransport(address NetAddress) *LocalTransport {
	return &LocalTransport{
		address:        address,
		consumeChannel: make(chan RPC, 1024),
		peers:          make(map[NetAddress]*LocalTransport),
	}
}

func (localTransport *LocalTransport) Consume() <-chan RPC {
	return localTransport.consumeChannel
}

// interface becuase mulitple options for transport eg local, tcp, udp, websocket
func (localTransport *LocalTransport) Connect(transport Transport) error {
	localTransport.lock.Lock()
	defer localTransport.lock.Unlock()

	localTransport.peers[transport.Address()] = transport.(*LocalTransport)

	return nil
}

func (localTransport *LocalTransport) SendMessage(to NetAddress, payload []byte) error {
	//multiple go routines can be read
	localTransport.lock.RLock()
	defer localTransport.lock.RUnlock()

	peer, ok := localTransport.peers[to]
	if !ok {
		return fmt.Errorf("%s coul not send message to unknown peers %s", localTransport.address, to)
	}

	peer.consumeChannel <- RPC{
		From:    localTransport.address,
		Payload: bytes.NewReader(payload),
	}

	return nil
}

func (ltransport *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range ltransport.peers {
		if err := ltransport.SendMessage(peer.Address(), payload); err != nil {
			return err
		}
	}
	return nil
}

func (ltransport *LocalTransport) Address() NetAddress {
	return ltransport.address
}
