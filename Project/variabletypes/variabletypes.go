package variabletypes

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

type ImAliveMsg struct {

}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type Elevator struct {
	Floor int
}

type Msg struct {
	//elev Elevator
	Messsage string
}