package fsm

import(
	."fmt"
	."time"
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
	elevatorObjectCh <-chan variableTypes.ElevatorObject,
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
			eventDoorTimeOut()

		case <- elevatorStuckTimer.C:
			eventElevatorStuckTimeOut()

		case SingleElevatorOrders := <-ordersCh:
			eventNewOrder()

		case SingleElevator.Floor := <-reachedFloorCh:
			eventNewFloor()
		}
	}
}

func eventArrivedAtFloor(){
	switch state {
	case MOVING:
		elevatorStuckTimer.Stop()
		if orderLogic.CheckForStop(singleElevator, singleElevatorOrders) {

		}
	}
}

func eventDoorTimeOut(removeOrderCh <-chan int){
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

func eventElevatorStuckTimeOut(){

}