package fsm

import(
	"fmt"
	"time"
	"elevio"
	"orderLogic"
	"../../config"
	"../../variableTypes"
)

type elevatorState int

const (
	IDLE elevatorState = iota
	OPEN
	MOVING
)

var state elevatorState
var singleElevator variableTypes.ElevatorObject
var singleElevatorOrders variableTypes.SingleOrderMatrix 

func Fsm(
	ordersCh <-chan variableTypes.SingleOrderMatrix,
	elevatorObjectCh chan<- variableTypes.ElevatorObject,
	removeOrderCh <-chan int,
	reachedFloorCh <-chan int
	) {

	state = IDLE
	//singleElevator.Floor := elevio.getFloor()
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()
	elevatorStuckTimer := time.NewTimer(3 * time.Second)
	elevatorStuckTimer.Stop()

	for {
		select {
		case <- doorTimer.C:
			fsmDoorTimeOut()

		case <- elevatorStuckTimer.C:
			fsmElevatorStuckTimeOut()

		case SingleElevatorOrders := <-ordersCh:
			fsmNewOrder()

		case SingleElevator.Floor := <-reachedFloorCh:
			fsmReachedFloor()
		}
	}
}

func fsmReachedFloor(){
	elevatorStuckTimer.Stop()
	switch state {
	case MOVING:
		if orderLogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetMotorDirection(variableTypes.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset()
			state = OPEN
		} else {
			elevatorStuckTimer.start()
		}
	}
}

func fsmDoorTimeOut(removeOrderCh <-chan int){
	switch state {
	case OPEN:
		elevio.SetDoorOpenLamp(false)
		removeOrderCh <- SingleElevator.Floor
		singleElevator.Dirn = orderLogic.chooseNextDirection(singleElevator, singleElevatorOrders)
		if singleElevator.dirn == variableTypes.MD_Stop {
			state = IDLE
			// elevatorStuckTimer.stop()
		} else {
			elevio.SetMotorDirection(singleElevator.dirn)
			elevatorStuckTimer.Reset(3 * time.Second)
			state = MOVING
		}
	}
}

func fsmElevatorStuckTimeOut(){
	fmt.Println("****************   ELEVATOR ENGINE ERROR!   ****************")
	fmt.Println("****************   RESTART ELEVATOR %d      ****************", config.ElevatorId)
	elevio.SetMotorDirection(variableTypes.MD_Stop)
	time.Sleep(time.Second * 1)
	//os.Exit(1)
	// Automatic restart? Other handling? how when what huh... todo
}