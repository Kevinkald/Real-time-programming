package utilities

import(
	"../../variabletypes"
	"../../config"
	"strconv"
	"fmt"
	//"time"
)

func CreateMapCopy(elevMap variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo {
	copyMap := make(variabletypes.AllElevatorInfo)
	for key, value := range elevMap {
		copyMap[key] = value
	}
	return copyMap
}

func InitMap()variabletypes.AllElevatorInfo{
	elevMap := make(map[string]variabletypes.SingleElevatorInfo)
	for id := 1; id <= config.N_Elevators; id++{
		id_string := strconv.Itoa(id)
		elevMap[id_string] = variabletypes.SingleElevatorInfo{}
	}
	return elevMap
}

func PrintMap(a variabletypes.AllElevatorInfo){
		for id := 1; id <= config.N_Elevators; id++{
			id_string := strconv.Itoa(id)
			fmt.Println("Elevator id: ",id_string)
			for floor := 0; floor < config.N_Floors; floor++{
				fmt.Println(a[id_string].OrderMatrix[floor])
			}
			fmt.Println("State", a[id_string].ElevObj.State)
			fmt.Println("Floor", a[id_string].ElevObj.Floor)
			fmt.Println("Dirn", a[id_string].ElevObj.Dirn)
		}
		//time.Sleep(200*time.Millisecond)
}