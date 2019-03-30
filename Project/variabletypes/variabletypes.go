package variabletypes

import(
	"../config"

)

type MotorDirection int

const (
	MDUp   MotorDirection = 1
	MDDown                = -1
	MDStop                = 0
)

type ButtonType int

const (
	BTHallUp   ButtonType = iota
	BTHallDown            
	BTCab                 
)

type ElevatorState int

const (
	IDLE ElevatorState = iota
	OPEN
	MOVING
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ButtonPress struct {
	floor int
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type ElevatorObject struct {
	Floor int
	Dirn MotorDirection
	State ElevatorState
}

type SingleOrderMatrix [config.NFloors][config.NButtons]bool

type SingleElevatorInfo struct {
	OrderMatrix SingleOrderMatrix
	ElevObj ElevatorObject
}

type AllElevatorInfo map[string]SingleElevatorInfo

type NetworkMsg struct{
	Id 	string
	Info AllElevatorInfo
}