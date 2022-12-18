package network

type NetAddress string

// module on server, needs access to all messages sent over transport layers
type Transport interface {
	// return channel of RPC
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddress, []byte) error
	Broadcast([]byte) error
	Address() NetAddress
}
