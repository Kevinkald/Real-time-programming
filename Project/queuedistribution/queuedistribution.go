package queuedistribution

import(
	"fmt"
	//"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"../fsm/elevio"
	//"../costfunction"
)

const (
	invalidIP string  = ""
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.AllElevatorInfo,
							NetworkMessageBroadcastCh chan<-  variabletypes.AllElevatorInfo,
							ButtonsCh <-chan variabletypes.ButtonEvent) {


	elevMap := make(map[string]variabletypes.SingleElevatorInfo)
	elevMap[config.ElevatorId] = variabletypes.SingleElevatorInfo{}

	//Send initialized elevMap to broadcasting
	//Important to copy the dynamic map before sending over channel
	msg := utilities.CreateMapCopy(elevMap)
	NetworkMessageBroadcastCh<- msg

	for {
		select{
		case p := <-peerUpdateCh:
			fmt.Println("Current alive nodes:",p.Peers)

		case b:= <-ButtonsCh:
		fmt.Println("Pushed button: {floor,type} ", b)
		elevio.SetButtonLamp(b.Button, b.Floor, true)

		case n := <-networkMessageCh:
		//Receive the maps from broadcast, do stuff with it here
		fmt.Println(n)
		}
	}
}


/*func DelegateOrder(elevMap variabletypes.AllElevatorInfo, button variabletypes.ButtonType, floor int) ElevatorId {
	AllElevMap := utilities.CreateMapCopy(elevMap)
	currentIP := invalidIP
	currentDuration := 0

	for id, info := range AllElevMap {
		currentElevator = AllElevatorInfo[i]
		elevDuration = costfunction.timeToServeRequest(currentElevator, button, floor)

		if elevDuration <= currentDuration {
			currentDuration = elevDuration
			currentIP = i

			currentElevator.OrderMatrix[floor][button] = 1
		}

	}

	return currentIP
}
*/

//if DelegateOrder = invalidIP 
//error


















