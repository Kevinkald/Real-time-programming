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
							removeOrderCh <-chan variabletypes.ButtonEvent) {


	elevMap := utilities.InitMap()

	//Just to check if elev obj syncs up
	var tmp = elevMap[config.ElevatorId]
	tmp.ElevObj.Floor = 3
	tmp.ElevObj.Dirn = 	variabletypes.MD_Up
	tmp.ElevObj.State = variabletypes.MOVING
	elevMap[config.ElevatorId] = tmp

	printMapCh := make(chan variabletypes.AllElevatorInfo)
	go utilities.PrintMap(printMapCh)


	//Send initialized elevMap to broadcasting
	//Important to copy the dynamic map before sending over channel
	msg := utilities.CreateMapCopy(elevMap)
	NetworkMessageBroadcastCh<- msg

	for {
		select{
			//WHY DOES THIS FLICKER WHEN PRINTING??
		//case p := <-peerUpdateCh:
			//fmt.Println("Current alive nodes:",p.Peers)

		case b:= <-ButtonsCh:
			fmt.Println("Pushed button: {floor,type} ", b)
			var tmp = elevMap[config.ElevatorId]
			tmp.OrderMatrix[b.Floor][b.Button] = true
			elevMap[config.ElevatorId] = tmp
			elevio.SetButtonLamp(b.Button, b.Floor, true)
			//Broadcast changes
			msg := utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg

		case n := <-networkMessageCh:
			//fmt.Println(n)
			elevMap = synchronizationlogic.Synchronize(elevMap,n)
			//Broadcast changes and print
			msg := utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
			time.Sleep(1*time.Millisecond)
			printMapCh <- msg
		
		case r := <-removeOrderCh:
			var tmp = elevMap[config.ElevatorId]
			tmp.OrderMatrix[r.Floor][r.Button] = false
			elevMap[config.ElevatorId] = tmp
			//Broadcast changes
			msg := utilities.CreateMapCopy(elevMap)
			NetworkMessageBroadcastCh<- msg
		}
	}
}