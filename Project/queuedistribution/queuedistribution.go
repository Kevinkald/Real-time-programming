package queuedistribution

import(
	"fmt"
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"../fsm/elevio"
	"./synchronizationlogic"
	"./costfunction"
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.NetworkMsg,
							NetworkMessageBroadcastCh chan<-  variabletypes.NetworkMsg,
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
	networkMessageTicker := time.NewTicker(time.Millisecond * 50)
	ticker := time.NewTicker(time.Millisecond * 500)


	//Send initialized elevMap to broadcasting
	//Important to copy the dynamic map before sending over channel
	var msg variabletypes.NetworkMsg
	var p variabletypes.PeerUpdate

	msg.Info = utilities.CreateMapCopy(elevMap)
	msg.Id = config.ElevatorId
	NetworkMessageBroadcastCh<- msg
	fmt.Println("Starting")

	for {
		select{
			//WHY DOES THIS FLICKER WHEN PRINTING??
		case new_p := <-peerUpdateCh: 
			p = new_p
			fmt.Println("Current alive nodes:",p.Peers)

		case b:= <-ButtonsCh:
			fmt.Println("Pushed button: {floor,type} ", b)
			

			// find best elevator to take order and set corresponding queue 
			chosenElevator := costfunction.DelegateOrder(elevMap, p, b, msg.Id)
			if chosenElevator == config.InvalidId {
				fmt.Println("Error: invalid Id")
			}
			var tmptwo = elevMap[chosenElevator]
			tmptwo.OrderMatrix[b.Floor][b.Button] = true
			elevMap[chosenElevator] = tmptwo

			// hva skjer her?
			var tmp = elevMap[config.ElevatorId]
			tmp.OrderMatrix[b.Floor][b.Button] = true
			elevMap[config.ElevatorId] = tmp

			elevio.SetButtonLamp(b.Button, b.Floor, true)

			//Broadcast changes
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
			if (len(p.Peers)== 0){
				ordersCh <- elevMap[config.ElevatorId].OrderMatrix
			}

		case n := <-networkMessageCh:
			//fmt.Println(n)
			elevMap = synchronizationlogic.Synchronize(elevMap,n.Info)
			//Broadcast changes and print
			//msg.Info := utilities.CreateMapCopy(elevMap)
			//NetworkMessageBroadcastCh<- msg
			//time.Sleep(1*time.Millisecond)

			ordersCh <- elevMap[config.ElevatorId].OrderMatrix
		
		case r := <-removeOrderCh:
			//todo: make this nicer
			var tmp = elevMap[config.ElevatorId]
			utilities.ClearOrder(r, tmp) 
			elevMap[config.ElevatorId] = tmp

			//Broadcast changes
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
		
		case q := <- elevatorObjectCh:
			var tmp = elevMap[config.ElevatorId]
			tmp.ElevObj = q
			elevMap[config.ElevatorId] = tmp

		case <- ticker.C:
			//utilities.PrintMap(utilities.CreateMapCopy(elevMap))

		case <-networkMessageTicker.C:
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
			time.Sleep(1*time.Millisecond)	
		}
	}
}


