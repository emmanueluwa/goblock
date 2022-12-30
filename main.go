package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/emmanueluwa/goblock/core"
	"github.com/emmanueluwa/goblock/crypto"
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

var transports = []network.Transport{
	network.NewLocalTransport("LOCAL"),
	// network.NewLocalTransport(("REMOTE_A")),
	// network.NewLocalTransport(("REMOTE_B")),
	// network.NewLocalTransport(("REMOT_C")),
	// ,
}

func main() {
	initRemoteServers(transports)
	localNode := transports[0]
	lateTransport := network.NewLocalTransport("LATE_NODE")
	// RemoteNodeA := transports[1]
	// RemoteNodeC := transports[3]

	// go func() {
	// 	for {
	// 		// transportRemote.SendMessage(transportLocal.Address(), []byte("Obavan people"))
	// 		if err := sendTransaction(RemoteNodeA, localNode.Address()); err != nil {
	// 			logrus.Error(err)
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	go func() {
		time.Sleep(7 * time.Second)
		lateServer := makeServer(string(lateTransport.Address()), lateTransport, nil)
		go lateServer.Start()
	}()

	//local server
	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", transports[0], &privKey)
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
		Transport:  transport,
		PrivateKey: pk,
		ID:         id,
		Transports: transports,
	}

	server, err := network.NewServer(options)
	if err != nil {
		log.Fatal(err)
	}
	return server
}

func sendGetStatusMessage(transport network.Transport, to network.NetAddress) error {
	var (
		getStatusMessage = new(network.GetStatusMessage)
		buffer           = new(bytes.Buffer)
	)
	if err := gob.NewDecoder(buffer).Encode(getStatusMessage); err != nil {
		return err
	}

	message := network.NewMessage(network.MessageTypeGetStatus, buffer.Bytes())
}

// placeholder for demonstration
func sendTransaction(transport network.Transport, to network.NetAddress) error {
	privKey := crypto.GeneratePrivateKey()
	transaction := core.NewTransaction(contract())
	transaction.Sign(privKey)
	buffer := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTxEncoder(buffer)); err != nil {
		return err
	}

	message := network.NewMessage(network.MessageTypeTx, buffer.Bytes())

	return transport.SendMessage(to, message.Bytes())
}

func contract() []byte {
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x03, 0x0a, 0x0d, 0xae}
	data = append(data, pushFoo...)

	return data
}
