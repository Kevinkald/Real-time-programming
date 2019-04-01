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
			synchedElevInfo[elevId] = 
				utilities.SetSingleElevatorObject(synchedElevInfo[elevId],receivedElevInfo[elevId].ElevObj)
		}

		// Synchronize orders
		for floor := 0; floor < config.NFloors; floor++{
			for button := 0; button < config.NButtons; button++{
				//If the setButtonLamp matrices have different values(true-false or false-true)
				if ((localElevInfo[elevId].OrderMatrix[floor][button])!=
					(receivedElevInfo[elevId].OrderMatrix[floor][button])){
					synchedElevInfo[elevId] = utilities.SetSingleElevatorMatrixValue(synchedElevInfo[elevId], floor,button, true);
					//If the local elev info is the one not having the order
					if (!localElevInfo[elevId].OrderMatrix[floor][button]){
						if ((localElevInfo[elevId].ElevObj.State==variabletypes.OPEN)&&
							(localElevInfo[elevId].ElevObj.Floor==floor)){
							synchedElevInfo[elevId] = utilities.SetSingleElevatorMatrixValue(synchedElevInfo[elevId], floor,button, false);
						}
					//If the received elev info is the one not having the order
					} else if((receivedElevInfo[elevId].ElevObj.State==variabletypes.OPEN)&&
							(receivedElevInfo[elevId].ElevObj.Floor==floor)){ 
						synchedElevInfo[elevId] = utilities.SetSingleElevatorMatrixValue(synchedElevInfo[elevId], floor,button, false);
					}
				}
			}
		}
	}
	return synchedElevInfo
}

func SynchronizeButtonLamps(	elevatorsCh <-chan variabletypes.AllElevatorInfo,
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
                    for button := 0; button < config.NButtons; button++{
                        setButtonLamp := false
                        for _,id := range peers.Peers {
                            if (elevators[id].OrderMatrix[floor][button]) {
                                if (button != int(variabletypes.BTCab))||(id == config.ElevatorId ) {
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