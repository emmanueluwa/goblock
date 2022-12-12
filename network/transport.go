package network

type NetAddress string

type RPC struct {
	From    NetAddress
	Payload []byte
}

// module on server, needs access to all messages sent over transport layers
type Transport interface {
	// return chanel of RPC
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddress, []byte) error
	Address() NetAddress
}
