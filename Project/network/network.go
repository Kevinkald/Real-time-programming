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
				NetworkMessageCh chan<-  variabletypes.AllElevatorInfo,
				NetworkMessageBroadcastCh <-chan  variabletypes.AllElevatorInfo) {

	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)

	//Start transmitting the elevator id to port
	go peers.Transmitter(config.PeerPort, config.ElevatorId, peerTxEnable)
	//Pass received network messages to peerUpdateCh
	go peers.Receiver(config.PeerPort, peerUpdateCh)

	//Start broadcasting messages received on NetworkMessageBroadcastCh
	go bcast.Transmitter(config.BroadcastPort, NetworkMessageBroadcastCh)
	//Pass received networkmessages to NetWorkMessageCh
	go bcast.Receiver(config.BroadcastPort, NetworkMessageCh)
}
// ikke send videre mld om det er din egen, dvs. endre mldtype