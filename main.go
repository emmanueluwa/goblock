package main

import (
	"bytes"
	"fmt"
	"log"
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
	transportRemoteA := network.NewLocalTransport(("REMOTE_A"))
	transportRemoteB := network.NewLocalTransport(("REMOTE_B"))
	transportRemoteC := network.NewLocalTransport(("REMOT_C"))

	transportLocal.Connect(transportRemoteA)
	transportRemoteA.Connect(transportRemoteB)
	transportRemoteB.Connect(transportRemoteC)
	transportRemoteA.Connect(transportLocal)

	initRemoteServers([]network.Transport{transportRemoteA, transportRemoteB, transportRemoteC})

	go func() {
		for {
			// transportRemote.SendMessage(transportLocal.Address(), []byte("Obavan people"))
			if err := sendTransaction(transportRemoteA, transportLocal.Address()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// go func() {
	// 	time.Sleep(7 * time.Second)
	// 	transportLate := network.NewLocalTransport("LATE_REMOTE")
	// 	transportRemoteC.Connect(transportLate)
	// 	lateServer := makeServer(string(transportLate.Address()), transportLate, nil)

	// 	go lateServer.Start()
	// }()

	//local server
	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", transportLocal, &privKey)
	localServer.Start()
}

func initRemoteServers(transports []network.Transport) {
	for i := 0; i < len(transports); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		server := makeServer(id, transports[i], nil)
		go server.Start()
	}
}

func makeServer(id string, transport network.Transport, pk *crypto.PrivateKey) *network.Server {
	options := network.ServerOptions{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{transport},
	}

	server, err := network.NewServer(options)
	if err != nil {
		log.Fatal(err)
	}
	return server
}

// placeholder for demonstration
func sendTransaction(transport network.Transport, to network.NetAddress) error {
	privKey := crypto.GeneratePrivateKey()
	// data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f} //, 0x03, 0x0a, 0x02, 0x0a, 0x0e}

	transaction := core.NewTransaction(data)
	transaction.Sign(privKey)
	buffer := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTxEncoder(buffer)); err != nil {
		return err
	}

	message := network.NewMessage(network.MessageTypeTx, buffer.Bytes())

	return transport.SendMessage(to, message.Bytes())
}
