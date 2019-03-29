package network

import (
	"./bcast"
	"./peers"
	"../variabletypes"
	"../config"
)

func Network(	peerUpdateCh chan<- variabletypes.PeerUpdate, 
				networkMessageCh chan<-  variabletypes.NetworkMsg,
				networkMessageBroadcastCh <-chan  variabletypes.NetworkMsg) {

	// We can disable/enable the transmitter after it has been started.
	peerTxEnable := make(chan bool)

	// Start transmitting and receiving peers
	go peers.Transmitter(config.PeerPort, config.ElevatorId, peerTxEnable)
	go peers.Receiver(config.PeerPort, peerUpdateCh)

	// Start transmitting and receiving elevator data
	go bcast.Transmitter(config.BroadcastPort, networkMessageBroadcastCh)
	go bcast.Receiver(config.BroadcastPort, networkMessageCh)
}