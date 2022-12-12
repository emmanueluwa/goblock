package main

import (
	"time"

	"github.com/emmanueluwa/goblock/network"
)

/***

DYNAMIC GENERIC BLOCKCHAIN
- CAN BE CONFIGURED
- IMPLEMENTATIONS CAN BE WRITTEN FOR TRANSPORT EG

SERVER(CONTAINER)
- MODULES, TRANSPORT LAYER -> TCP. UDP
           BLOCK
					 TRANSACTION
					 KEYPAIRS

***/

func main() {
	transportLocal := network.NewLocalTransport("LOCAL")
	transportRemote := network.NewLocalTransport(("REMOTE"))

	transportLocal.Connect(transportRemote)
	transportRemote.Connect(transportLocal)

	go func() {
		for {
			transportRemote.SendMessage(transportLocal.Address(), []byte("Obavan people"))
			time.Sleep(1 * time.Second)
		}
	}()

	//configure our own server/node
	options := network.ServerOptions{
		Transports: []network.Transport{transportLocal},
	}

	server := network.NewServer(options)
	server.Start()
}
