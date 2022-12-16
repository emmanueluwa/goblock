package network

import (
	"crypto"
	"fmt"
	"time"

	"github.com/emmanueluwa/goblock/core"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type ServerOptions struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOptions
	//for server to know when to create block from mempool values
	memPool     *TxPool
	isValidator bool
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(options ServerOptions) *Server {

	if options.BlockTime == time.Duration(0) {
		options.BlockTime = defaultBlockTime
	}
	if options.RPCDecodeFunc == nil {
		options.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	server := &Server{
		ServerOptions: options,
		memPool:       NewTxPool(),
		isValidator:   options.PrivateKey != nil,
		rpcChannel:    make(chan RPC),
		quitChannel:   make(chan struct{}, 1),
	}

	//if no rpc processer is given from options, the server is default processor
	if server.RPCProcessor == nil {
		server.RPCProcessor = server
	}

	return server
}

func (server *Server) Start() {
	server.initTransports()
	ticker := time.NewTicker(server.BlockTime)

free:
	for {
		//block until a case can run and is executed
		select {
		case rpc := <-server.rpcChannel:
			message, err := server.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}

			if err := server.RPCProcessor.ProcessMessage(message); err != nil {
				logrus.Error(err)
			}
		case <-server.quitChannel:
			break free
		case <-ticker.C:
			if server.isValidator {
				server.createNewBlock()
			}

		}
	}
	fmt.Println("Server shutdown")
}

func (server *Server) ProcessMessage(message *DecodedMessage) error {
	switch t := message.Data.(type) {
	case *core.Transaction:
		return server.ProcessTransaction(t)
	}
	return nil
}

func (server *Server) ProcessTransaction(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})

	if server.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")

		return nil
	}

	if err := transaction.Verify(); err != nil {
		return err
	}

	transaction.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           transaction.Hash(core.TxHasher{}),
		"mempool length": server.memPool.Len(),
	}).Info("adding new transaction to mempool")

	//TODO (@emmanueluwa): broadcast this transaction to peers in network

	return server.memPool.Add(transaction)
}

// if validator check mempool, put all transactions into a block
func (server *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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
