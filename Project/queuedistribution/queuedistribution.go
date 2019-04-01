package queuedistribution

import(
	"fmt"
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"./synchlogic"
	"./costfunction"
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.NetworkMsg,
							networkMessageBroadcastCh chan<-  variabletypes.NetworkMsg,
							buttonsCh <-chan variabletypes.ButtonEvent,
							removeOrderCh <-chan int,
							ordersCh chan<- variabletypes.SingleOrderMatrix,
		 					elevatorObjectCh <-chan variabletypes.ElevatorObject,
		 					elevatorsCh chan<- variabletypes.AllElevatorInfo,
		 					alivePeersCh chan<- variabletypes.PeerUpdate) {

	elevatorMap := utilities.InitMap()
	networkMessageTicker := time.NewTicker(time.Millisecond * 15)
	orderChannelTicker := time.NewTicker(time.Millisecond * 100)

	var msg variabletypes.NetworkMsg
	var peers variabletypes.PeerUpdate

	msg.Info = utilities.CreateMapCopy(elevatorMap)
	msg.Id = config.ElevatorId

	networkMessageBroadcastCh<- msg

	for {
		select{
		case p := <-peerUpdateCh: 
			receivedPeers := p
			if (len(receivedPeers.Peers)!=len(peers.Peers)){
				redistributed_orders := redistributeOrders(receivedPeers,elevatorMap)
				elevatorMap = redistributed_orders
			}
			peers = receivedPeers
			alivePeersCh <- peers

		case b:= <-buttonsCh:
			chosenElevatorID := costfunction.DelegateOrder(elevatorMap, peers, b)
			if chosenElevatorID == config.InvalidId {
				fmt.Println("Error: Invalid Id")
			}

			var chosenElevator = elevatorMap[chosenElevatorID]
			chosenElevator.OrderMatrix[b.Floor][b.Button] = true
			elevatorMap[chosenElevatorID] = chosenElevator

			msg.Info = utilities.CreateMapCopy(elevatorMap)
			networkMessageBroadcastCh<- msg

		case n := <-networkMessageCh:
			elevatorMap = synchlogic.SynchronizeElevInfo(elevatorMap,n.Info)
		
		case r := <-removeOrderCh:
			var elevator = elevatorMap[config.ElevatorId]

			for button := 0; button < config.NButtons; button++{
				elevator.OrderMatrix[r][button] = false
			}
			elevatorMap[config.ElevatorId] = elevator

			msg.Info = utilities.CreateMapCopy(elevatorMap)
			networkMessageBroadcastCh<- msg
		
		case q := <-elevatorObjectCh:
			var elevator = elevatorMap[config.ElevatorId]
			elevator.ElevObj = q
			elevatorMap[config.ElevatorId] = elevator

		case <-networkMessageTicker.C:
			msg.Info = utilities.CreateMapCopy(elevatorMap)
			networkMessageBroadcastCh<- msg
			time.Sleep(1*time.Millisecond)	

		case <-orderChannelTicker.C:
			if (elevatorMap[config.ElevatorId].ElevObj.State != variabletypes.OPEN){
				ordersCh <- elevatorMap[config.ElevatorId].OrderMatrix
			}
			elevators := utilities.CreateMapCopy(elevatorMap)
			elevatorsCh<- elevators
		}
	}
}

func redistributeOrders( peers variabletypes.PeerUpdate,
						 elevatorMap variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo {
	
	redistributedMap := utilities.CreateMapCopy(elevatorMap)
	var redistributedOrder variabletypes.ButtonEvent

	for _,lostElevatorId := range peers.Lost {
		for floor := 0; floor < config.NFloors; floor++ {
			redistributedOrder.Floor = floor
			for button := variabletypes.BTHallUp; button <= variabletypes.BTHallDown; button++{
				redistributedOrder.Button = button
				
				if (elevatorMap[lostElevatorId].OrderMatrix[floor][button]){
					new_id := costfunction.DelegateOrder(elevatorMap, peers, redistributedOrder)
					elevator := redistributedMap[new_id]
					elevator.OrderMatrix[floor][button] = true
					redistributedMap[new_id] = elevator
				}
			}
		}
	}
	return redistributedMap
}