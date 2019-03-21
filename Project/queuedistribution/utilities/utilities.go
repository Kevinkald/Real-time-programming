package utilities

import(
	"../../variabletypes"
	"../../config"
	"strconv"
	"fmt"
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

func PrintMap(elevMap variabletypes.AllElevatorInfo){
	for id := 1; id <= config.N_Elevators; id++{
		id_string := strconv.Itoa(id)
		fmt.Println("Elevator id: ",id_string)
		for floor := 0; floor < config.M_Floors; floor++{
			fmt.Println(elevMap[id_string].OrderMatrix[floor])
		}
		fmt.Println("State", elevMap[id_string].ElevObj.State)
		fmt.Println("Floor", elevMap[id_string].ElevObj.Floor)
		fmt.Println("Dirn", elevMap[id_string].ElevObj.Dirn)
	}
}