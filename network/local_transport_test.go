package network

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(test *testing.T) {
	localTransportA := NewLocalTransport("A")
	localTransportB := NewLocalTransport("B")

	localTransportA.Connect(localTransportB)
	localTransportB.Connect(localTransportA)
	assert.Equal(test, localTransportA.peers[localTransportB.address], localTransportB)
	assert.Equal(test, localTransportB.peers[localTransportA.address], localTransportA)
}

func TestSendMessage(test *testing.T) {
	localTransportA := NewLocalTransport("A")
	localTransportB := NewLocalTransport("B")

	localTransportA.Connect(localTransportB)
	localTransportB.Connect(localTransportA)

	message := []byte("Obavan people")
	assert.Nil(test, localTransportA.SendMessage(localTransportB.address, message))

	//consume channel, test real scenario of outside world having no access to consumer
	rpc := <-localTransportB.Consume()
	r, err := ioutil.ReadAll(rpc.Payload)

	assert.Nil(test, err)
	assert.Equal(test, r, message)
	//check address
	assert.Equal(test, rpc.From, localTransportA.address)
}

func TestBroadcast(test *testing.T) {
	transportA := NewLocalTransport("A")
	transportB := NewLocalTransport("B")
	transportC := NewLocalTransport("C")

	transportA.Connect(transportB)
	transportA.Connect(transportC)

	message := []byte("meow")
	assert.Nil(test, transportA.Broadcast(message))

	rpcB := <-transportB.Consume()
	b, err := ioutil.ReadAll(rpcB.Payload)
	assert.Nil(test, err)
	assert.Equal(test, b, message)

	rpcC := <-transportC.Consume()
	c, err := ioutil.ReadAll(rpcC.Payload)
	assert.Nil(test, err)
	assert.Equal(test, c, message)
}
