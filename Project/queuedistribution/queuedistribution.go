package queuedistribution

import(
	"fmt"
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.AllElevatorInfo,
							NetworkMessageBroadcastCh chan<-  variabletypes.AllElevatorInfo,
							/*ButtonsCh <-chan variabletypes.ButtonEvent*/) {

	fmt.Println(config.K_Buttons)

	elevMap := make(map[string]variabletypes.SingleElevatorInfo)
	elevMap[config.ElevatorId] = variabletypes.SingleElevatorInfo{}
	

	for {	
		//Important to copy the dynamic map before sending over channel
		msg := utilities.CreateMapCopy(elevMap)
		NetworkMessageBroadcastCh<- msg
		time.Sleep(1 * time.Second)
	}
}