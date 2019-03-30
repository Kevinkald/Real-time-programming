package synchlogic

import(
	"../../variabletypes"
	"../../config"
	"../utilities"
	"../../fsm/elevio"
	"time"
)

func SynchronizeElevInfo(	localElevInfo variabletypes.AllElevatorInfo,
							receivedElevInfo variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo{

	synchedElevInfo := utilities.CreateMapCopy(localElevInfo)

	for elevId, _ := range localElevInfo{

		// Synchronize elevator objects
		if (elevId != config.ElevatorId){
			var tmp = synchedElevInfo[elevId]
			tmp.ElevObj = receivedElevInfo[elevId].ElevObj
			synchedElevInfo[elevId] = tmp
		}

		// Synchronize orders
		for floor := 0; floor < config.NFloors; floor++{
			for button := 0; button < config.NButtons; button++{
				//If the setButtonLamp matrices have different values(true-false or false-true)
				if ((localElevInfo[elevId].OrderMatrix[floor][button])!=
					(receivedElevInfo[elevId].OrderMatrix[floor][button])){
					//Set setButtonLamp to true(union)
					var tmp = synchedElevInfo[elevId]
					tmp.OrderMatrix[floor][button] = true
					synchedElevInfo[elevId] = tmp 
					//If the local elev info is the one not having an setButtonLamp
					if (!localElevInfo[elevId].OrderMatrix[floor][button]){
						if ((localElevInfo[elevId].ElevObj.State==variabletypes.OPEN)&&
							(localElevInfo[elevId].ElevObj.Floor==floor)){
							var tmp = synchedElevInfo[elevId]
							tmp.OrderMatrix[floor][button] = false
							synchedElevInfo[elevId] = tmp 
						}
					//If the received elev info is the one not having an setButtonLamp
					} else if((receivedElevInfo[elevId].ElevObj.State==variabletypes.OPEN)&&
							(receivedElevInfo[elevId].ElevObj.Floor==floor)){
						var tmp = synchedElevInfo[elevId]
						tmp.OrderMatrix[floor][button] = false
						synchedElevInfo[elevId] = tmp 
					}
				}
			}
		}
	}
	return synchedElevInfo
}

func SynchronizeButtonLamps(elevatorsCh <-chan variabletypes.AllElevatorInfo,
							alivePeersCh <-chan variabletypes.PeerUpdate){

    var peers variabletypes.PeerUpdate
    var elevators variabletypes.AllElevatorInfo
    ticker := time.NewTicker(time.Millisecond * 100)

    for {
        select {
            case e := <-elevatorsCh:
            	elevators = e
            case <-ticker.C:
                for floor := 0; floor < config.NFloors; floor++{
                    for button := variabletypes.BTHallUp; button <= variabletypes.BTCab; button++{
                        setButtonLamp := false
                        for _,id := range peers.Peers {
                            if (elevators[id].OrderMatrix[floor][button]) {
                                if (button != variabletypes.BTCab)||(id == config.ElevatorId ) {
                                    setButtonLamp = true
                                }
                            }
                        }
                        elevio.SetButtonLamp(variabletypes.ButtonType(button), floor, setButtonLamp)
                    }
                }
            case p := <-alivePeersCh:
            	peers = p
        }     
    }
}