package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/emmanueluwa/goblock/core"
)

/***
Recieved plain byte payload needs to be converted into
	some kind of message, so we know what we need to do with it
***/

// *No enums in GO but similar,
type MessageType byte

const (
	MessageTypeTx MessageType = 0x1
	//automatically increments since it is below the above, similar to enum
	MessageTypeBlock
)

//*

type RPC struct {
	From    NetAddress
	Payload io.Reader
}

type Message struct {
	//in header first bytes allows us to know the type of data
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

// our own message/protocol
func (message *Message) Bytes() []byte {
	buffer := &bytes.Buffer{}
	gob.NewEncoder(buffer).Encode(message)
	return buffer.Bytes()
}

type DecodedMessage struct {
	From NetAddress
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	message := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	//find out message type so we can handle it
	switch message.Header {
	case MessageTypeTx:
		transaction := new(core.Transaction)
		if err := transaction.Decode(core.NewGobTxDecoder(bytes.NewReader(message.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: transaction,
		}, nil

	// dealing with messageType that we do not accept for example
	default:
		return nil, fmt.Errorf("invalid message header %x", message.Header)
	}
}

// process encoded data from handler
type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}
