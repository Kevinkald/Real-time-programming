package queuedistribution

import(
	"fmt"
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"./synchlogic"
	"./costfunction"
	"../orderlogic"
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
				redistributed_orders := orderlogic.RedistributeOrders(receivedPeers,elevatorMap)
				elevatorMap = redistributed_orders
			}
			peers = receivedPeers
			alivePeersCh <- peers

		case b:= <-buttonsCh:
			chosenElevatorID := costfunction.DelegateOrder(elevatorMap, peers, b)
			if chosenElevatorID == config.InvalidId {
				fmt.Println("Error: Invalid Id")
			}

			elevatorMap[chosenElevatorID] = 
			utilities.SetSingleElevatorMatrixValue(elevatorMap[chosenElevatorID], int(b.Floor), int(b.Button), true);

			msg.Info = utilities.CreateMapCopy(elevatorMap)
			networkMessageBroadcastCh<- msg

		case n := <-networkMessageCh:
			elevatorMap = synchlogic.SynchronizeElevInfo(elevatorMap,n.Info)
		
		case r := <-removeOrderCh:
			for button := 0; button < config.NButtons; button++{
				elevatorMap[config.ElevatorId] = 
				utilities.SetSingleElevatorMatrixValue(elevatorMap[config.ElevatorId], r, button, false);
			}

			msg.Info = utilities.CreateMapCopy(elevatorMap)
			networkMessageBroadcastCh<- msg
		
		case q := <-elevatorObjectCh:
			var elevator = elevatorMap[config.ElevatorId]
			elevator.ElevObj = q
			elevatorMap[config.ElevatorId] = elevator

		case <-networkMessageTicker.C:
			msg.Info = utilities.CreateMapCopy(elevatorMap)
			networkMessageBroadcastCh<- msg

		case <-orderChannelTicker.C:
			if (elevatorMap[config.ElevatorId].ElevObj.State != variabletypes.OPEN){
				ordersCh <- elevatorMap[config.ElevatorId].OrderMatrix
			}
			elevators := utilities.CreateMapCopy(elevatorMap)
			elevatorsCh<- elevators
		}
	}
}

func RedistributeOrders( peers variabletypes.PeerUpdate,
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
					redistributedMap[new_id] = 
					utilities.SetSingleElevatorMatrixValue(redistributedMap[new_id], floor, int(button), true);
				}
			}
		}
	}
	return redistributedMap
}