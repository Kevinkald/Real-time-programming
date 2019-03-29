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

	fmt.Println("Fsm goroutines up and running")

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
			fmt.Println("received orders(FSM)")
			elevatorObjectCh <- singleElevator
			fmt.Println("sent ElevObj (FSM)")

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
			if singleElevator.Dirn != variabletypes.MDStop{							// BAD FIX!!!! Fix this with more clear communication
				elevio.SetMotorDirection(singleElevator.Dirn)							// between modules
				elevatorStuckTimerResetCh <- true
				singleElevator.State = variabletypes.MOVING
				fmt.Println("MOVING!!!!!!!!!!!!!!!!!!!!!!!!!")
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
		fmt.Println("Closing doors")
		for btn := 0; btn < config.NButtons; btn++{
			singleElevatorOrders[singleElevator.Floor][btn] = false
		}
		removeOrderCh <- singleElevator.Floor
		fmt.Println("Order removed")
		elevio.SetDoorOpenLamp(false)
		singleElevator.Dirn = orderlogic.ChooseNextDirection(singleElevator, singleElevatorOrders)
		if singleElevator.Dirn == variabletypes.MDStop {
			singleElevator.State = variabletypes.IDLE
			fmt.Println("Elevator in IDLE")
			elevatorStuckTimerStopCh <- true
		} else {
			elevio.SetMotorDirection(singleElevator.Dirn)
			elevatorStuckTimerResetCh <- true
			singleElevator.State = variabletypes.MOVING
		}
	}
}

func fsmElevatorStuckTimeOut(){
	fmt.Println("****************   ELEVATOR ENGINE ERROR!   ****************")
	fmt.Println("****************   RESTART ELEVATOR %d      ****************", config.ElevatorId)
	elevio.SetMotorDirection(variabletypes.MDStop)
	//time.Sleep(time.Second * 1)
	(exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go "+config.ElevatorId+" "+config.ElevatorPort)).Run()
	os.Exit(1)
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