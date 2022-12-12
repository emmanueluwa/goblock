package network

import (
	"fmt"
	"time"
)

type ServerOptions struct {
	Transports []Transport
}

type Server struct {
	ServerOptions

	rpcChannel chan RPC

	quitChannel chan struct{}
}

func NewServer(options ServerOptions) *Server {
	return &Server{
		ServerOptions: options,
		rpcChannel:    make(chan RPC),
		quitChannel:   make(chan struct{}, 1),
	}
}

func (server *Server) Start() {
	server.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		//block until a case can run and is executed
		select {
		case rpc := <-server.rpcChannel:
			fmt.Printf("%+v\n", rpc)
		case <-server.quitChannel:
			break free
		case <-ticker.C:
			fmt.Println("do stuff ever x seconds")
		}
	}
	fmt.Println("Server shutdown")
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
