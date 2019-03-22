package fsm

import(
	"fmt"
	"time"
	"./elevio"
	"../config"
	"../orderlogic"
	"../variabletypes"
)

var state variabletypes.ElevatorState
var singleElevator variabletypes.ElevatorObject
var singleElevatorOrders variabletypes.SingleOrderMatrix 

func Fsm(ordersCh <-chan variabletypes.SingleOrderMatrix, elevatorObjectCh chan<- variabletypes.ElevatorObject, removeOrderCh chan<- variabletypes.ElevatorObject) {
	elevio.Init(config.HardwarePort)
	fmt.Println("Elevator initiated")
	state = variabletypes.IDLE
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

	fmt.Println("Goroutines up and running")

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
	fmt.Println("New order received!")
	switch state {
	case variabletypes.IDLE:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetDoorOpenLamp(true)
			doorTimerResetCh <- true
			state = variabletypes.OPEN
		} else {
			singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimerResetCh <- true
			state = variabletypes.MOVING
		}
	}
}

func fsmReachedFloor(doorTimerResetCh chan<- bool, elevatorStuckTimerResetCh chan<- bool, elevatorStuckTimerStopCh chan<- bool){
	elevatorStuckTimerStopCh <- true
	switch state {
	case variabletypes.MOVING:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetMotorDirection(variabletypes.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			doorTimerResetCh <- true
			state = variabletypes.OPEN
		} else {
			elevatorStuckTimerResetCh <- true
		}
	}
}

func fsmDoorTimeOut(removeOrderCh chan<- variabletypes.ElevatorObject, elevatorStuckTimerResetCh chan<- bool){
	switch state {
	case variabletypes.OPEN:
		fmt.Println("Closing doors")
		removeOrderCh <- singleElevator
		elevio.SetDoorOpenLamp(false)
		singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
		if singleElevator.Dirn == variabletypes.MD_Stop {
			state = variabletypes.IDLE
			fmt.Println("Elevator in IDLE")
			// elevatorStuckTimer.stop()
		} else {
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimerResetCh <- true
			state = variabletypes.MOVING
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