package utilities

import(
	"../../variabletypes"
)

func CreateMapCopy(elevMap variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo {
	copyMap := make(variabletypes.AllElevatorInfo)
	for key, value := range elevMap {
		copyMap[key] = value
	}
	return copyMap
}