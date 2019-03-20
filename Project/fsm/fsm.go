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

func Fsm(ordersCh <-chan variabletypes.SingleOrderMatrix, elevatorObjectCh chan<- variabletypes.ElevatorObject, removeOrderCh <-chan int) {

	state = IDLE
	//singleElevator.Floor := elevio.getFloor()
	doorTimer := time.NewTimer(2 * time.Second)
	doorTimer.Stop()

	elevatorStuckTimer := time.NewTimer(5 * time.Second)
	elevatorStuckTimer.Stop()
	reachedFloorCh := make(chan int)

	go elevio.PollFloorSensor(reachedFloorCh)

	for {
		select {
		case <- doorTimer.C:
			fsmDoorTimeOut(removeOrderCh)

		case <- elevatorStuckTimer.C:
			fsmElevatorStuckTimeOut()

		case msg1 := <-ordersCh:
			singleElevatorOrders = msg1
			fsmNewOrder()

		case msg2 := <-reachedFloorCh:
			singleElevator.Floor = msg2
			fsmReachedFloor()
		}
	}
}

func fsmNewOrder() {
	switch state {
	case OPEN:
		// TODO: What happens if an order is recieved in the same floor while the door is open?

	case IDLE:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset(2 * time.Second)
			state = OPEN
		} else {
			singleElevator.Dirn = orderlogic.chooseNextDirection(singleElevator, singleElevatorOrders)
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimer.Reset(5 * time.Second)
			state = MOVING
		}
	}
}

func fsmReachedFloor(){
	elevatorStuckTimer.Stop()
	switch state {
	case MOVING:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetMotorDirection(variabletypes.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset(2 * time.Second)
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
		singleElevator.Dirn = orderlogic.chooseNextDirection(singleElevator, singleElevatorOrders)
		if singleElevator.Dirn == variabletypes.MD_Stop {
			state = IDLE
			// elevatorStuckTimer.stop()
		} else {
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimer.Reset(5 * time.Second)
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
	// Automatic restart? Other handling? how when what huh... todo
}

func fsmDoorTimer(doorTimerResetCh <-chan bool, doorTimerOutCh chan<- bool){
	doorTimer := time.NewTimer(2 * time.Second)
	doorTimer.Stop()
	for{
		select{
		case <-doorTimerResetCh:
			timer.Stop()
			timer.Reset(2 * time.Second)
			case <-doorTimer.C:
				

		}
	}
}