package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/emmanueluwa/goblock/core"
	"github.com/sirupsen/logrus"
)

/***
Recieved plain byte payload needs to be converted into
	some kind of message, so we know what we need to do with it
***/

// *No enums in GO but similar,
type MessageType byte

const (
	//automatically increments since it is below the above, similar to enum
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3
	MessageTypeStatus    MessageType = 0x4
	MessageTypeGetStatus MessageType = 0x5
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

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": message.Header,
	}).Debug("new incoming message")

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

	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecoder(bytes.NewReader(message.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil

	case MessageTypeGetStatus:
		return &DecodedMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil

	case MessageTypeStatus:
		statusMessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(message.Data)).Decode(statusMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: statusMessage,
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
