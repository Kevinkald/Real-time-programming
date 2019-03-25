package queuedistribution

import(
	"fmt"
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"../fsm/elevio"
	"./synchronizationlogic"
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.AllElevatorInfo,
							NetworkMessageBroadcastCh chan<-  variabletypes.AllElevatorInfo,
							ButtonsCh <-chan variabletypes.ButtonEvent,
							removeOrderCh <-chan int,
							ordersCh chan<- variabletypes.SingleOrderMatrix,
		 					elevatorObjectCh <-chan variabletypes.ElevatorObject) {


	elevMap := utilities.InitMap()

	/*Just to check if elev obj syncs up
	var tmp = elevMap[config.ElevatorId]
	tmp.ElevObj.Floor = 2
	tmp.ElevObj.Dirn = 	variabletypes.MD_Up
	tmp.ElevObj.State = variabletypes.MOVING
	elevMap[config.ElevatorId] = tmp
	*/
	ticker := time.NewTicker(time.Millisecond * 500)


	//Send initialized elevMap to broadcasting
	//Important to copy the dynamic map before sending over channel
	msg := variabletypes.NetworkMsg
	msg.Info = utilities.CreateMapCopy(elevMap)
	msg.Id = config.ElevatorId
	NetworkMessageBroadcastCh<- msg
	fmt.Println("Starting")
	for {
		select{
			//WHY DOES THIS FLICKER WHEN PRINTING??
		case p := <-peerUpdateCh:
			fmt.Println("Current alive nodes:",p.Peers)

		case b:= <-ButtonsCh:
			fmt.Println("Pushed button: {floor,type} ", b)
			var tmp = elevMap[config.ElevatorId]

			// find best elevator to take order and set queue 
			chosenElevator := costfunction.DelegateOrder(elevMap, b)
			elevMap[chosenElevator].OrderMatrix[b.Floor][b.Button] = true

			tmp.OrderMatrix[b.Floor][b.Button] = true
			elevMap[config.ElevatorId] = tmp

			elevio.SetButtonLamp(b.Button, b.Floor, true)

			//Broadcast changes
			msg := utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
			if (len(p.Peers)== 0){
				ordersCh <- elevMap[config.ElevatorId].OrderMatrix
			}

		case n := <-networkMessageCh:
			//fmt.Println(n)
			elevMap = synchronizationlogic.Synchronize(elevMap,n)
			//Broadcast changes and print
			msg := utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
			time.Sleep(1*time.Millisecond)
		
		case r := <-removeOrderCh:
			//todo: make this nicer
			var tmp = elevMap[config.ElevatorId]
			utilities.ClearOrder(r, tmp.OrderMatrix) 
			elevMap[config.ElevatorId] = tmp

			//Broadcast changes
			msg := utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
		
		case q := <- elevatorObjectCh:
			var tmp = elevMap[config.ElevatorId]
			tmp.ElevObj = q
			elevMap[config.ElevatorId] = tmp
		case <- ticker.C:
			utilities.PrintMap(utilities.CreateMapCopy(elevMap))
		}
	}
}


