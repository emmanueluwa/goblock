package main

import (
	"bytes"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/emmanueluwa/goblock/core"
	"github.com/emmanueluwa/goblock/crypto"
	"github.com/emmanueluwa/goblock/network"
	"github.com/sirupsen/logrus"
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
			// transportRemote.SendMessage(transportLocal.Address(), []byte("Obavan people"))
			if err := sendTransaction(transportRemote, transportLocal.Address()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	privKey := crypto.GeneratePrivateKey()
	//configure our own server/node
	options := network.ServerOptions{
		PrivateKey: &privKey,
		ID:         "LOCAL",
		Transports: []network.Transport{transportLocal},
	}

	server, err := network.NewServer(options)
	if err != nil {
		return log.Fatal(err)
	}
	server.Start()
}

// placeholder for demonstration
func sendTransaction(transport network.Transport, to network.NetAddress) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(10000000)), 10))
	transaction := core.NewTransaction(data)
	transaction.Sign(privKey)
	buffer := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTxEncoder(buffer)); err != nil {
		return err
	}

	message := network.NewMessage(network.MessageTypeTx, buffer.Bytes())

	return transport.SendMessage(to, message.Bytes())
}
