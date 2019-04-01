package utilities

import(
	"../../variabletypes"
	"../../config"
	"strconv"
	"fmt"
)

func CreateMapCopy(elevatorMap variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo {
	copyMap := make(variabletypes.AllElevatorInfo)
	for key, value := range elevatorMap {
		copyMap[key] = value
	}
	return copyMap
}

func SetSingleElevatorMatrixValue(	elevatorMap variabletypes.SingleElevatorInfo,
									 floor int, button int, value bool)variabletypes.SingleElevatorInfo{
	elevatorMap.OrderMatrix[floor][button] = value
	return elevatorMap
}

func InitMap() variabletypes.AllElevatorInfo {
	elevatorMap := make(map[string]variabletypes.SingleElevatorInfo)
	for id := 1; id <= config.NElevators; id++ {
		id_string := strconv.Itoa(id)
		elevatorMap[id_string] = variabletypes.SingleElevatorInfo{}
	}
	return elevatorMap
}

func PrintMap(elevatorMap variabletypes.AllElevatorInfo){
		for id := 1; id <= config.NElevators; id++{
			id_string := strconv.Itoa(id)
			fmt.Println("Elevator id: ",id_string)
			for floor := 0; floor < config.NFloors; floor++{
				fmt.Println(elevatorMap[id_string].OrderMatrix[floor])
			}
			fmt.Println("State", elevatorMap[id_string].ElevObj.State)
			fmt.Println("Floor", elevatorMap[id_string].ElevObj.Floor)
			fmt.Println("Dirn", elevatorMap[id_string].ElevObj.Dirn)
		}
}