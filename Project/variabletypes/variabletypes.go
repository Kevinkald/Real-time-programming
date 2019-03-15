package variabletypes

import(
	"../config"

)

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
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
}

type SingleOrderMatrix [config.M_Floors][config.K_Buttons]bool

type SingleElevatorInfo struct {
	OrderMatrix SingleOrderMatrix
	State ElevatorObject
}

type AllElevatorInfo map[string]SingleElevatorInfo