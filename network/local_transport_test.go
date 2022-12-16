package network

import (
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
	buffer := make([]byte, len(message))
	n, err := rpc.Payload.Read(buffer)
	assert.Nil(test, err)
	assert.Equal(test, n, len(message))

	assert.Equal(test, buffer, message)
	//check address
	assert.Equal(test, rpc.From, localTransportA.address)
}
