package fsm

import(
	"fmt"
	"time"
	"./elevio"
	"../config"
	"../orderlogic"
	"../variabletypes"
	"os"
	"os/exec"
)

var singleElevator variabletypes.ElevatorObject
var singleElevatorOrders variabletypes.SingleOrderMatrix 

func Fsm(	ordersCh <-chan variabletypes.SingleOrderMatrix,
		 	elevatorObjectCh chan<- variabletypes.ElevatorObject,
		 	removeOrderCh chan<- int) {
	
	singleElevator.State = variabletypes.IDLE
	singleElevator.Dirn = variabletypes.MDStop

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
			fsmDoorTimeOut(removeOrderCh, elevatorStuckTimerResetCh, elevatorStuckTimerStopCh)
			elevatorObjectCh <- singleElevator

		case <- elevatorStuckTimerOutCh:
			fsmElevatorStuckTimeOut()

		case msg1 := <-ordersCh:
			singleElevatorOrders = msg1
			fsmNewOrder(doorTimerResetCh, elevatorStuckTimerResetCh)
			elevatorObjectCh <- singleElevator

		case msg2 := <-reachedFloorCh:
			singleElevator.Floor = msg2
			fsmReachedFloor(doorTimerResetCh, elevatorStuckTimerResetCh, elevatorStuckTimerStopCh)
			elevatorObjectCh <- singleElevator
		}
	}
}

func fsmNewOrder(doorTimerResetCh chan<- bool, elevatorStuckTimerResetCh chan<- bool) {
	switch singleElevator.State {
	case variabletypes.IDLE:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetDoorOpenLamp(true)
			doorTimerResetCh <- true
			singleElevator.State = variabletypes.OPEN
		} else {
			singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
			if singleElevator.Dirn != variabletypes.MDStop{
				elevio.SetMotorDirection(singleElevator.Dirn)		
				elevatorStuckTimerResetCh <- true
				singleElevator.State = variabletypes.MOVING
			}
		}
	}
}

func fsmReachedFloor(doorTimerResetCh chan<- bool, elevatorStuckTimerResetCh chan<- bool, elevatorStuckTimerStopCh chan<- bool){
	elevatorStuckTimerStopCh <- true
	switch singleElevator.State {
	case variabletypes.MOVING:
		if orderlogic.CheckForStop(singleElevator, singleElevatorOrders) {
			elevio.SetMotorDirection(variabletypes.MDStop)
			elevio.SetDoorOpenLamp(true)
			if (singleElevator.Floor == config.NFloors - 1 || singleElevator.Floor == 0){
				singleElevator.Dirn = variabletypes.MDStop
			}
			doorTimerResetCh <- true
			singleElevator.State = variabletypes.OPEN
		} else {
			elevatorStuckTimerResetCh <- true
		}
	}
}

func fsmDoorTimeOut(removeOrderCh chan<- int, elevatorStuckTimerResetCh chan<- bool, elevatorStuckTimerStopCh chan<- bool){
	switch singleElevator.State {
	case variabletypes.OPEN:
		for button := 0; button < config.NButtons; button++{
			singleElevatorOrders[singleElevator.Floor][button] = false
		}
		removeOrderCh <- singleElevator.Floor
		elevio.SetDoorOpenLamp(false)
		singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
		if singleElevator.Dirn == variabletypes.MDStop {
			elevatorStuckTimerStopCh <- true
			singleElevator.State = variabletypes.IDLE
		} else {
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimerResetCh <- true
			singleElevator.State = variabletypes.MOVING
		}
	}
}

func fsmElevatorStuckTimeOut(){
	fmt.Println("****************   ELEVATOR ENGINE ERROR. RESTARTING!   ****************")
	elevio.SetMotorDirection(variabletypes.MDStop)
	(exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go "+config.ElevatorId+" "+config.ElevatorPort)).Run()
	os.Exit(1)
}

func fsmDoorTimer(tesetCh <-chan bool, timerOutCh chan<- bool){
	doorTimer := time.NewTimer(config.DOOR_OPEN_TIME * time.Second)
	doorTimer.Stop()
	for{
		select{
		case <-resetCh:
			doorTimer.Reset(config.DOOR_OPEN_TIME * time.Second)
		case <-doorTimer.C:
			timerOutCh <- true
		}
	}
}

func fsmElevatorStuckTimer(resetCh <-chan bool, stopCh <-chan bool, timerOutCh chan<- bool){
	elevatorStuckTimer := time.NewTimer(5 * time.Second)
	elevatorStuckTimer.Stop()
	for {
		select{
		case <- stopCh:
			elevatorStuckTimer.Stop()
		case <- resetCh:
			elevatorStuckTimer.Reset(5 * time.Second)
		case <- elevatorStuckTimer.C:
			timerOutCh <- true
		}
	}
}