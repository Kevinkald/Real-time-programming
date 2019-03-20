package fsm

import(
	"fmt"
	"time"
	"./elevio"
	"../config"
	"../orderlogic"
	"../variabletypes"
)

type elevatorState int

const (
	IDLE elevatorState = iota
	OPEN
	MOVING
)

var state elevatorState
var singleElevator variabletypes.ElevatorObject
var singleElevatorOrders variabletypes.SingleOrderMatrix 

func Fsm(ordersCh <-chan variabletypes.SingleOrderMatrix, elevatorObjectCh chan<- variabletypes.ElevatorObject, removeOrderCh chan<- int) {

	state = IDLE
	//singleElevator.Floor := elevio.getFloor()
	elevatorStuckTimerResetCh := make(chan bool)
	elevatorStuckTimerStopCh := make(chan bool)
	elevatorStuckTimerOutCh := make(chan bool)

	doorTimerResetCh := make(chan bool)
	doorTimerOutCh := make (chan bool)

	reachedFloorCh := make(chan int)

	go fsmElevatorStuckTimer(elevatorStuckTimerResetCh, elevatorStuckTimerStopCh, elevatorStuckTimerOutCh)
	go fsmDoorTimer(doorTimerResetCh, doorTimerOutCh)
	go elevio.PollFloorSensor(reachedFloorCh)

	for {
		select {
		case <- doorTimerOutCh:
			fsmDoorTimeOut(removeOrderCh, elevatorStuckTimerResetCh)

		case <- elevatorStuckTimerOutCh:
			fsmElevatorStuckTimeOut()

		case msg1 := <-ordersCh:
			singleElevatorOrders = msg1
			fsmNewOrder(doorTimerResetCh, elevatorStuckTimerResetCh)

		case msg2 := <-reachedFloorCh:
			singleElevator.Floor = msg2
			fsmReachedFloor(doorTimerResetCh, elevatorStuckTimerResetCh, elevatorStuckTimerStopCh)
		}
	}
}

func fsmNewOrder(doorTimerResetCh chan<- bool, elevatorStuckTimerResetCh chan<- bool) {
	switch state {
	case IDLE:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetDoorOpenLamp(true)
			doorTimerResetCh <- true
			state = OPEN
		} else {
			singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimerResetCh <- true
			state = MOVING
		}
	}
}

func fsmReachedFloor(doorTimerResetCh chan<- bool, elevatorStuckTimerResetCh chan<- bool, elevatorStuckTimerStopCh chan<- bool){
	elevatorStuckTimerStopCh <- true
	switch state {
	case MOVING:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetMotorDirection(variabletypes.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			doorTimerResetCh <- true
			state = OPEN
		} else {
			elevatorStuckTimerResetCh <- true
		}
	}
}

func fsmDoorTimeOut(removeOrderCh chan<- int, elevatorStuckTimerResetCh chan<- bool){
	switch state {
	case OPEN:
		elevio.SetDoorOpenLamp(false)
		removeOrderCh <- singleElevator.Floor
		singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
		if singleElevator.Dirn == variabletypes.MD_Stop {
			state = IDLE
			// elevatorStuckTimer.stop()
		} else {
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimerResetCh <- true
			state = MOVING
		}
	}
}

func fsmElevatorStuckTimeOut(){
	fmt.Println("****************   ELEVATOR ENGINE ERROR!   ****************")
	fmt.Println("****************   RESTART ELEVATOR %d      ****************", config.ElevatorId)
	elevio.SetMotorDirection(variabletypes.MD_Stop)
	time.Sleep(time.Second * 1)
	//os.Exit(1)
}

func fsmDoorTimer(doorTimerResetCh <-chan bool, doorTimerOutCh chan<- bool){
	doorTimer := time.NewTimer(2 * time.Second)
	doorTimer.Stop()
	for{
		select{
		case <-doorTimerResetCh:
			doorTimer.Stop()
			doorTimer.Reset(2 * time.Second)
		case <-doorTimer.C:
			doorTimerOutCh <- true
		}
	}
}

func fsmElevatorStuckTimer(elevatorStuckTimerResetCh <-chan bool, elevatorStuckTimerStopCh <-chan bool, elevatorStuckTimerOutCh chan<- bool){
	elevatorStuckTimer := time.NewTimer(5 * time.Second)
	elevatorStuckTimer.Stop()
	for {
		select{
		case <- elevatorStuckTimerStopCh:
			elevatorStuckTimer.Stop()
		case <- elevatorStuckTimerResetCh:
			elevatorStuckTimer.Stop()
			elevatorStuckTimer.Reset(5 * time.Second)
		case <- elevatorStuckTimer.C:
			elevatorStuckTimerOutCh <- true
		}
	}
}