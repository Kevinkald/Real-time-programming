package fsm

import(
	"time"
	"../config"
	"../variabletypes"
)

func ElevatorStuckTimer(resetCh <-chan bool, stopCh <-chan bool, timerOutCh chan<- bool){
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

func DoorTimer(tesetCh <-chan bool, timerOutCh chan<- bool){
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