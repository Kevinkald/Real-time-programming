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

func SetSingleElevatorObject(   elevatorMap variabletypes.SingleElevatorInfo,
                                updatedObject variabletypes.ElevatorObject)variabletypes.SingleElevatorInfo{
    elevatorMap.ElevObj = updatedObject
    return elevatorMap
}

func InitMap() variabletypes.AllElevatorInfo {
	elevatorMap := make(map[string]variabletypes.SingleElevatorInfo)
	for id := 1; id <= config.NElevators; id++ {
		stringId := strconv.Itoa(id)
		elevatorMap[stringId] = variabletypes.SingleElevatorInfo{}
	}
	return elevatorMap
}

func PrintMap(elevatorMap variabletypes.AllElevatorInfo){
		for id := 1; id <= config.NElevators; id++{
			stringId := strconv.Itoa(id)
			fmt.Println("Elevator id: ",stringId)
			for floor := 0; floor < config.NFloors; floor++{
				fmt.Println(elevatorMap[stringId].OrderMatrix[floor])
			}
			fmt.Println("State", elevatorMap[stringId].ElevObj.State)
			fmt.Println("Floor", elevatorMap[stringId].ElevObj.Floor)
			fmt.Println("Dirn", elevatorMap[stringId].ElevObj.Dirn)
		}
}