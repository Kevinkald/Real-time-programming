package synchlogic

import(
	"../../variabletypes"
	"../../config"
	"../utilities"
	"fmt"
	"../../fsm/elevio"
	"time"
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
		for floor := 0; floor < config.N_Floors; floor++{

			for button := 0; button < config.N_Buttons; button++{

				//If the two queues have different values(true-false,false-true)
				if (e_local[elevid].OrderMatrix[floor][button]!=e_received[elevid].OrderMatrix[floor][button]){
					fmt.Println("synch logic: entered true-false cond")
					//Set queue element to true(union)
					var tmp = e_synched[elevid]
					tmp.OrderMatrix[floor][button] = true
					e_synched[elevid] = tmp 
					//If the local is the one having false
					if (!e_local[elevid].OrderMatrix[floor][button]){
						fmt.Println("synch logic: entered local one having the false")
						//If local elev is also OPEN, remove order if in that corresponding floor
						if ((e_local[elevid].ElevObj.State==variabletypes.OPEN)&&
							(e_local[elevid].ElevObj.Floor==floor)){
							fmt.Println("synch logic: remove order")
							var tmp = e_synched[elevid]
							tmp.OrderMatrix[floor][button] = false
							e_synched[elevid] = tmp 
							//elevio.SetButtonLamp(variabletypes.ButtonType(button), floor, false)
						}
							//elevio.SetButtonLamp(variabletypes.ButtonType(button), floor, true)
						//If the received is the one having false
					} else if((e_received[elevid].ElevObj.State==variabletypes.OPEN)&&
							(e_received[elevid].ElevObj.Floor==floor)){
						fmt.Println("synch logic: entered received one having the false")
						fmt.Println("synch logic: remove order")
						var tmp = e_synched[elevid]
						tmp.OrderMatrix[floor][button] = false
						e_synched[elevid] = tmp 
						//elevio.SetButtonLamp(variabletypes.ButtonType(button), floor, false)
					}
				}
			}
		}
	}
	return e_synched
}

func SynchronizeButtonLamps(elevatorsCh <-chan variabletypes.AllElevatorInfo,
							alivePeersCh <-chan variabletypes.PeerUpdate){
    var peers variabletypes.PeerUpdate
    for {
        select {
            case elevators := <-elevatorsCh:
                for floor := 0; floor < config.N_Floors; floor++{ // for all floors
                    for btn := variabletypes.BT_HallUp; btn <= variabletypes.BT_Cab; btn++{ // and all buttons
                        order := false
                        for _,id := range peers.Peers { // for all alive elevators
                            if (elevators[id].OrderMatrix[floor][btn]) { // if there is an order
                               	//order = true
                                if (btn != variabletypes.BT_Cab)|| (id == config.ElevatorId ) { // if it's not a cabcall or it is your id
                                    order = true
                                }
                            }
                        }
                        if (order) {
                            elevio.SetButtonLamp(variabletypes.ButtonType(btn), floor, true)
                        } else {
                            elevio.SetButtonLamp(variabletypes.ButtonType(btn), floor, false)
                            }
                    }
                }
                time.Sleep(20*time.Millisecond)
            case p := <-alivePeersCh:
            	peers = p
        }
            
    }
}