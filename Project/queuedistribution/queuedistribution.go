package queuedistribution

import(
	"fmt"
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
	//"../fsm/elevio"
	"./synchlogic"
	"./costfunction"
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.NetworkMsg,
							NetworkMessageBroadcastCh chan<-  variabletypes.NetworkMsg,
							ButtonsCh <-chan variabletypes.ButtonEvent,
							removeOrderCh <-chan int,
							ordersCh chan<- variabletypes.SingleOrderMatrix,
		 					elevatorObjectCh <-chan variabletypes.ElevatorObject,
		 					elevatorsCh chan<- variabletypes.AllElevatorInfo) {


	elevMap := utilities.InitMap()

	/*Just to check if elev obj syncs up
	var tmp = elevMap[config.ElevatorId]
	tmp.ElevObj.Floor = 2
	tmp.ElevObj.Dirn = 	variabletypes.MD_Up
	tmp.ElevObj.State = variabletypes.MOVING
	elevMap[config.ElevatorId] = tmp
	*/
	ticker := time.NewTicker(time.Millisecond * 1000)
	networkMessageTicker := time.NewTicker(time.Millisecond * 15)
	orderChannelTicker := time.NewTicker(time.Millisecond * 100)


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
		case new_p := <-peerUpdateCh: 
			p = new_p
			fmt.Println("Current alive nodes:",p.Peers)

		case b:= <-ButtonsCh:
			fmt.Println("Pushed button: {floor,type} ", b)
			

			// find best elevator to take order and set corresponding queue 
			chosenElevator := costfunction.DelegateOrder(elevMap, p, b)


			if chosenElevator == config.InvalidId {
				fmt.Println("Error: invalid Id, order lost")
			}


			var tmp = elevMap[chosenElevator]
			tmp.OrderMatrix[b.Floor][b.Button] = true
			elevMap[chosenElevator] = tmp

			//elevio.SetButtonLamp(b.Button, b.Floor, true)

			//Broadcast changes
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg

			
			//if (len(p.Peers)== 0){
			/*
			fmt.Println("button sending")
			if (elevMap[config.ElevatorId].ElevObj.State != variabletypes.OPEN){
				ordersCh <- elevMap[config.ElevatorId].OrderMatrix
			}
			fmt.Println("button sent")*/
			//}


		case n := <-networkMessageCh:

			//fmt.Println(n)
			elevMap = synchlogic.Synchronize(elevMap,n.Info)

			//Broadcast changes and print
			//msg.Info := utilities.CreateMapCopy(elevMap)
			//NetworkMessageBroadcastCh<- msg
			//time.Sleep(1*time.Millisecond)
			/*
			fmt.Println("net sending")
			if (elevMap[config.ElevatorId].ElevObj.State != variabletypes.OPEN){
				ordersCh <- elevMap[config.ElevatorId].OrderMatrix
			}
			fmt.Println("net sent")
			*/
			//Only synch if the received message is not sent by this node
			//if (n.Id!=config.ElevatorId){
				//ordersCh <- elevMap[config.ElevatorId].OrderMatrix
			//}
			
			//Broadcast changes and print 						NB: This should be done after given time intervals

			//msg.Info = utilities.CreateMapCopy(elevMap)
			//NetworkMessageBroadcastCh<- msg
			//time.Sleep(1*time.Millisecond)

			//Make nicer!
		
		case r := <-removeOrderCh:
			//todo: make this nicer
			var tmp = elevMap[config.ElevatorId]

			for button := 0; button < config.N_Buttons; button++{
				tmp.OrderMatrix[r][button] = false
				//elevio.SetButtonLamp(variabletypes.ButtonType(button), r, false)
			}
			elevMap[config.ElevatorId] = tmp

			//Broadcast changes
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
		
		case q := <-elevatorObjectCh:
			var tmp = elevMap[config.ElevatorId]
			tmp.ElevObj = q
			elevMap[config.ElevatorId] = tmp

		case <-ticker.C:
			utilities.PrintMap(utilities.CreateMapCopy(elevMap))

		case <-networkMessageTicker.C:
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
			time.Sleep(1*time.Millisecond)	
		case <-orderChannelTicker.C:
			if (elevMap[config.ElevatorId].ElevObj.State != variabletypes.OPEN){
				ordersCh <- elevMap[config.ElevatorId].OrderMatrix
			}
			elevators := utilities.CreateMapCopy(elevMap)
			elevatorsCh<- elevators
		}
	}
}