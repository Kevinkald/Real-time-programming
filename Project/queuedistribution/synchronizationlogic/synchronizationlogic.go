import (
	"../../variabletypes"
	"../../config"
)

func synchronize(	e_local variabletypes.AllElevatorInfo,
					e_received variabletypes.AllElevatorInfo)
					variabletypes.AllElevatorInfo)variabletypes.AllElevatorInfo{

	e_synched := e_local

	//Loop through all local elevators
	for elevid, local_elevinfo := range e_local{

		//1. Synchronize elevator object by updating
		if (elevid != config.ElevatorId){
			e_synched[elevid].Elevobj = e_received[elevid].Elevobj//This will prob fail
		}

		//2. Synchronize queues
		//Loop through all elements in queues
		for floor := 0; floor < config.M_Floors; floor++{

			for button := 0; button < config.K_Buttons; button++{


				if ()

			}
		}

	}

	return e_synched
}