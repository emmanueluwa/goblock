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

type RPCHandler interface {
	//decoder
	HandleRPC(rpc RPC) error
}

type DefaultRPCHandler struct {
	processor RPCProcessor
}

func NewDefaultRPCHandler(processor RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{
		processor: processor,
	}
}

// our implementation, therefor no need to abstract way to decode it
func (handler *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	message := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&message); err != nil {
		return fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	//find out message type so we can handle it
	switch message.Header {
	case MessageTypeTx:
		transaction := new(core.Transaction)
		if err := transaction.Decode(core.NewGobTxDecoder(bytes.NewReader(message.Data))); err != nil {
			return err
		}
		return handler.processor.ProcessTransaction(rpc.From, transaction)
	// dealing with messageType that we do not accept for example
	default:
		return fmt.Errorf("invalid message header %x", message.Header)
	}
}

// process encoded data from handler
type RPCProcessor interface {
	ProcessTransaction(NetAddress, *core.Transaction) error
}
