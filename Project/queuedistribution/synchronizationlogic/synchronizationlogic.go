package synchronizationlogic

import(
	"../../variabletypes"
	"../../config"
	"../utilities"
	//"fmt"
)

func Synchronize(	e_local variabletypes.AllElevatorInfo,
					e_received variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo{

	e_synched := utilities.CreateMapCopy(e_local)

	//Loop through all elev id's
	for elevid, _ := range e_local{

		//1. Synchronize elevator objects
		if (elevid != config.ElevatorId){

			var tmp = e_synched[elevid]
			tmp.ElevObj = e_received[elevid].ElevObj
			e_synched[elevid] = tmp
		}

		//2. Synchronize queues
		//Loop through all elements in queues
		for floor := 0; floor < config.M_Floors; floor++{

			for button := 0; button < config.K_Buttons; button++{

				//If the two queues have different values(true-false,false-true)
				if (e_local[elevid].OrderMatrix[floor][button]!=e_received[elevid].OrderMatrix[floor][button]){
					//Set queue element to true(union)
					var tmp = e_synched[elevid]
					tmp.OrderMatrix[floor][button] = true
					e_synched[elevid] = tmp 
					//If the local is the one having false
					if (!e_local[elevid].OrderMatrix[floor][button]){
						//If local elev is also IDLE||OPEN, remove order if in that corresponding floor
						if ((e_local[elevid].ElevObj.State==variabletypes.IDLE)||
							(e_local[elevid].ElevObj.State==variabletypes.OPEN)){
							if (e_local[elevid].ElevObj.Floor==floor){
								var tmp = e_synched[elevid]
								tmp.OrderMatrix[floor][button] = false
								e_synched[elevid] = tmp 
							}
						}
					} else if((e_received[elevid].ElevObj.State==variabletypes.IDLE))||
							(e_received[elevid].ElevObj.State==variabletypes.OPEN){
							if (e_received[elevid].ElevObj.Floor==floor){
								var tmp = e_synched[elevid]
								tmp.OrderMatrix[floor][button] = false
								e_synched[elevid] = tmp 
							}
					}
				}
			}
		}
	}

	return e_synched
}