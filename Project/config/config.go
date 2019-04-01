package config

import (
	"os"
	"time"
)

var ElevatorPort string
var ElevatorId string

func ConfigInit(){
	ElevatorId = os.Args[1]
	ElevatorPort = os.Args[2]
}

const (
	PeerUpdateInterval = 15*time.Millisecond
	PeerTimeout = 1000*time.Millisecond

	NFloors int = 4
	NButtons int = 3
	NElevators int = 3

	PeerPort int = 17563
	BroadcastPort int = 17564

	DoorOpenTime = 2*time.Second
	StuckTime = 5*time.Second
	PollRate = 20*time.Millisecond

	BroadcastInterval = time.Millisecond * 15
	OrderUpdateInterval = time.Millisecond * 100
)