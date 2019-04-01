package fsm

import(
	"time"
	"../config"
)

func ElevatorStuckTimer(resetCh <-chan bool, stopCh <-chan bool, timerOutCh chan<- bool){
	elevatorStuckTimer := time.NewTimer(config.StuckTime)
	elevatorStuckTimer.Stop()
	for {
		select{
		case <- stopCh:
			elevatorStuckTimer.Stop()
		case <- resetCh:
			elevatorStuckTimer.Reset(config.StuckTime)
		case <- elevatorStuckTimer.C:
			timerOutCh <- true
		}
	}
}

func DoorTimer(resetCh <-chan bool, timerOutCh chan<- bool){
	doorTimer := time.NewTimer(config.DoorOpenTime)
	doorTimer.Stop()
	for{
		select{
		case <-resetCh:
			doorTimer.Reset(config.DoorOpenTime)
		case <-doorTimer.C:
			timerOutCh <- true
		}
	}
}