package queuedistribution

import(
	"fmt"
	//"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"../fsm/elevio"
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