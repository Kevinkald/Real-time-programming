package network

import (
	"./bcast"
	//"./localip"
	"./peers"
	"../variabletypes"
	//"flag"
	//"fmt"
	//"os"
	//"time"
	//"strings"
	"../config"
)

func Network(	peerUpdateCh chan<- variabletypes.PeerUpdate, 
				NetworkMessageCh chan<-  variabletypes.NetworkMessage,
				NetworkMessageBroadcastCh <-chan  variabletypes.NetworkMessage) {

	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)

	//Start transmitting the elevator id to port
	go peers.Transmitter(15647, config.ELEVATOR_ID, peerTxEnable)
	//Pass received network messages to peerUpdateCh
	go peers.Receiver(15647, peerUpdateCh)

	//Start broadcasting messages received on NetworkMessageBroadcastCh
	go bcast.Transmitter(16569, NetworkMessageBroadcastCh)
	//Pass received networkmessages to NetWorkMessageCh
	go bcast.Receiver(16569, NetworkMessageCh)
}

	/*go func() {
		Msg := variabletypes.NetworkMessage{}
		for {
			msgTx <- Msg
			time.Sleep(1 * time.Second)
		}
	}()*/

	//

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	/*if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		config.ELEVATOR_ID = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}*/