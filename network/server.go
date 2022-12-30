package network

import (
	"bytes"
	"crypto"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/emmanueluwa/goblock/core"
	"github.com/emmanueluwa/goblock/types"
	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type ServerOptions struct {
	ID            string
	Transport     Transport
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOptions
	//for server to know when to create block from mempool values
	memPool     *TransactionPool
	chain       *core.Blockchain
	isValidator bool
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(options ServerOptions) (*Server, error) {

	if options.BlockTime == time.Duration(0) {
		options.BlockTime = defaultBlockTime
	}
	if options.RPCDecodeFunc == nil {
		options.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if options.Logger == nil {
		options.Logger = log.NewLogfmtLogger(os.Stderr)
		options.Logger = log.With(options.Logger, "address", options.Transport.Address())
	}

	chain, err := core.NewBlockchain(options.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	server := &Server{
		ServerOptions: options,
		chain:         chain,
		memPool:       NewTransactionPool(1000),
		isValidator:   options.PrivateKey != nil,
		rpcChannel:    make(chan RPC),
		quitChannel:   make(chan struct{}, 1),
	}

	//if no rpc processer is given from options, the server is default processor
	if server.RPCProcessor == nil {
		server.RPCProcessor = server
	}

	if server.isValidator {
		go server.validatorLoop()
	}

	server.bootstrapNodes()

	return server, nil
}

func (server *Server) Start() {
	server.initTransports()

free:
	for {
		//block until a case can run and is executed
		select {
		case rpc := <-server.rpcChannel:
			message, err := server.RPCDecodeFunc(rpc)
			if err != nil {
				server.Logger.Log("error", err)
			}

			if err := server.RPCProcessor.ProcessMessage(message); err != nil {
				if err == core.ErrBlockKnown {
					server.Logger.Log("error", err)
				}
			}
		case <-server.quitChannel:
			break free
		}
	}
	server.Logger.Log("message", "Server is shutting down")
}

func (server *Server) bootstrapNodes() {
	for _, transport := range server.Transports {
		if server.Transport.Address() != transport.Address() {
			if err := server.Transport.Connect(transport); err != nil {
				server.Logger.Log("error", "could not connect to remote", "err", err)
			}
			server.Logger.Log("message", "connect to remote", "our address", server.Transport.Address(), "address", transport.Address())

			//send getStatusMessage for sync if needed
			fmt.Printf("%s sending message to => %+s\n", server.Transport.Address(), transport.Address())

			if err := server.sendGetStatusMessage(transport); err != nil {
				server.Logger.Log("error", "sendGetStatusMessage", "err", err)
			}
		}
	}
}

func (server *Server) validatorLoop() {
	ticker := time.NewTicker(server.BlockTime)

	server.Logger.Log("message", "Starting validator loop", "blockTime", server.BlockTime)

	for {
		<-ticker.C
		server.createNewBlock()
	}
}

func (server *Server) ProcessMessage(message *DecodedMessage) error {
	fmt.Printf("recieving message: %+v\n", message.Data)
	switch t := message.Data.(type) {
	case *core.Transaction:
		return server.processTransaction(t)
	case *core.Block:
		return server.processBlock(t)
	case *GetStatusMessage:
		return server.processGetStatusMessage(message.From, t)
	case *StatusMessage:
		return server.processStatusMessage(message.From, t)
	case *GetBlocksMessage:
		return server.processGetBlocksMessage(message.From, t)
	}

	return nil
}

func (server *Server) processGetBlocksMessage(from NetAddress, data *GetBlocksMessage) error {
	fmt.Printf("recieved get blocks message %+v", data)

	return nil
}

// TODO(@emmanueluwa): REMOVE LOGIC FROM MAIN FUNCTION TO HERE
// TRANSPORT THAT IS OUR OWN SHOUDL SUFFICE
func (server *Server) sendGetStatusMessage(transport Transport) error {
	var (
		getStatusMessage = new(GetStatusMessage)
		buffer           = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buffer).Encode(getStatusMessage); err != nil {
		return err
	}

	message := NewMessage(MessageTypeGetStatus, buffer.Bytes())
	if err := server.Transport.SendMessage(transport.Address(), message.Bytes()); err != nil {
		return err
	}

	return nil
}

func (server *Server) broadcast(payload []byte) error {
	for _, transport := range server.Transports {
		if err := transport.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) processStatusMessage(from NetAddress, data *StatusMessage) error {
	//only ask for blocks if block height(of node we recieve statusMessage from) is higher
	if data.CurrentHeight <= server.chain.Height() {
		server.Logger.Log("message", "cannot sync blockHeight too low", "ourHeight", server.chain.Height(), "address", from)
		return nil
	}

	//for case of node having blocks higher than us
	getBlocksMessage := &GetBlocksMessage{
		From: server.chain.Height(),
		To:   0,
	}

	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(getBlocksMessage); err != nil {
		return err
	}

	message := NewMessage(MessageTypeGetBlocks, buffer.Bytes())

	return server.Transport.SendMessage(from, message.Bytes())
}

func (server *Server) processGetStatusMessage(from NetAddress, data *GetStatusMessage) error {
	fmt.Printf("received GETstatus message from %s => %+v\n", from, data)

	statusMessage := &StatusMessage{
		CurrentHeight: server.chain.Height(),
		ID:            server.ID,
	}

	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(statusMessage); err != nil {
		return err
	}

	message := NewMessage(MessageTypeStatus, buffer.Bytes())

	return server.Transport.SendMessage(from, message.Bytes())
}

func (server *Server) processBlock(block *core.Block) error {
	if err := server.chain.AddBlock(block); err != nil {
		return err
	}

	go server.broadcastBlock(block)

	return nil
}

func (server *Server) processTransaction(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})

	if server.memPool.Contains(hash) {
		return nil
	}

	if err := transaction.Verify(); err != nil {
		return err
	}

	server.Logger.Log(
		"message", "adding new transaction to mempool",
		"hash", hash,
		"mempoolLength", server.memPool.PendingCount(),
	)

	go server.broadcastTransaction(transaction)

	server.memPool.Add(transaction)

	return nil
}

func (server *Server) broadcastBlock(block *core.Block) error {
	buffer := &bytes.Buffer{}
	if err := block.Encode(core.NewGobBlockEncoder(buffer)); err != nil {
		return err
	}

	message := NewMessage(MessageTypeBlock, buffer.Bytes())

	return server.broadcast(message.Bytes())
}

// encoding broadcast for transport
func (server *Server) broadcastTransaction(transaction *core.Transaction) error {
	buffer := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTxEncoder(buffer)); err != nil {
		return err
	}

	message := NewMessage(MessageTypeTx, buffer.Bytes())

	return server.broadcast(message.Bytes())
}

func (server *Server) initTransports() {
	// for every transport we have, make them listen for messages
	// spin up a go routine and keep reading, consuming
	// pipe every rpc from go routine, from each transport, directly into server and own rpcChannel
	for _, transport := range server.Transports {
		go func(transport Transport) {
			for rpc := range transport.Consume() {
				server.rpcChannel <- rpc
			}
		}(transport)
	}
}

// if validator check mempool, put all transactions into a block
func (server *Server) createNewBlock() error {
	currentHeader, err := server.chain.GetHeader(server.chain.Height())
	if err != nil {
		return err
	}

	//using all transactions in mempool for now till internal structure of transaction
	// is known, complexity function will be implemented to determine no. of transactions
	// to include in each block
	transactions := server.memPool.Pending()

	block, err := core.NewBlockFromPrevHeader(currentHeader, transactions)
	if err != nil {
		return err
	}

	if err := block.Sign(*server.PrivateKey); err != nil {
		return err
	}

	if err := server.chain.AddBlock(block); err != nil {
		return err
	}

	server.memPool.ClearPending()

	go server.broadcastBlock(block)

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		TimeStamp: 000000,
	}

	block, _ := core.NewBlock(header, nil)
	return block
}
