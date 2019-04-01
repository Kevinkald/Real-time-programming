package queuedistribution

import(
	"time"
	"../config"
	"../variabletypes"
	"./utilities"
	"./synchlogic"
	"./orderassignment"
)

func Queuedistribution(		peerUpdateCh <-chan variabletypes.PeerUpdate,
							networkMessageCh <-chan variabletypes.NetworkMsg,
							networkMessageBroadcastCh chan<- variabletypes.NetworkMsg,
							buttonsCh <-chan variabletypes.ButtonEvent,
							removeOrderCh <-chan int,
							ordersCh chan<- variabletypes.SingleOrderMatrix,
		 					elevatorObjectCh <-chan variabletypes.ElevatorObject,
		 					elevatorsCh chan<- variabletypes.AllElevatorInfo,
		 					alivePeersCh chan<- variabletypes.PeerUpdate) {

	elevatorMap := utilities.InitMap()
	broadcastTicker := time.NewTicker(config.BroadcastInterval)
	orderChannelTicker := time.NewTicker(config.OrderUpdateInterval)

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
				redistributedOrders := orderassignment.RedistributeOrders(receivedPeers,elevatorMap)
				elevatorMap = redistributedOrders
			}
			peers = receivedPeers
			alivePeersCh <- peers

		case b:= <-buttonsCh:
			chosenElevatorID := orderassignment.DelegateOrder(elevatorMap, peers, b)

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

			elevatorMap[config.ElevatorId] = 
				utilities.SetSingleElevatorObject(elevatorMap[config.ElevatorId],q)

		case <-broadcastTicker.C:
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