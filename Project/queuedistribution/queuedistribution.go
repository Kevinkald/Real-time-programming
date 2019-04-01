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
		 					elevatorsCh chan<- variabletypes.AllElevatorInfo,
		 					alivePeersCh chan<- variabletypes.PeerUpdate) {

	elevMap := utilities.InitMap()

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

	for {
		select{
		case new_p := <-peerUpdateCh: 
			received_p := new_p
			if (len(received_p.Peers)!=len(p.Peers)){
				redistributed_orders := redistributeOrders(received_p,elevMap)
				elevMap = redistributed_orders
			}
			p = received_p
			alivePeersCh <- p

		case b:= <-ButtonsCh:
			// find best elevator to take order and set corresponding queue 
			chosenElevator := costfunction.DelegateOrder(elevMap, p, b)

			if chosenElevator == config.InvalidId {
				fmt.Println("Error: Invalid Id")
			}
			var tmp = elevMap[chosenElevator]
			tmp.OrderMatrix[b.Floor][b.Button] = true
			elevMap[chosenElevator] = tmp

			//Broadcast changes
			msg.Info = utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg

		case n := <-networkMessageCh:
			elevMap = synchlogic.SynchronizeElevInfo(elevMap,n.Info)
		
		case r := <-removeOrderCh:
			var tmp = elevMap[config.ElevatorId]

			for button := 0; button < config.NButtons; button++{
				tmp.OrderMatrix[r][button] = false
			}
			elevMap[config.ElevatorId] = tmp

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

func redistributeOrders( peers variabletypes.PeerUpdate,
						 elevMap variabletypes.AllElevatorInfo)variabletypes.AllElevatorInfo{
	redistMap := utilities.CreateMapCopy(elevMap)
	var redistributedOrder variabletypes.ButtonEvent
	for _,lostElevatorId := range peers.Lost {
		for floor := 0; floor < config.NFloors; floor++{
			redistributedOrder.Floor = floor
			for btn := variabletypes.BTHallUp; btn <= variabletypes.BTHallDown; btn++{
				redistributedOrder.Button = btn
				if (elevMap[lostElevatorId].OrderMatrix[floor][btn]){
					new_id := costfunction.DelegateOrder(elevMap, peers, redistributedOrder)
					tmp := redistMap[new_id]
					tmp.OrderMatrix[floor][btn] = true
					redistMap[new_id] = tmp
				}
			}
		}
	}
	return redistMap
}