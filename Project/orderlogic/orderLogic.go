package orderlogic

import(
	"../config"
	"../variabletypes"
)

func ordersAbove(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)bool{
	for floor := elevator.Floor + 1; floor < config.N_Floors; floor++ {
		for btn := 0; btn < config.N_Buttons; btn++ {
			if orders[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)bool{
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < config.N_Buttons; btn++ {
			if orders[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ChooseNextDirection(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)variabletypes.MotorDirection{
	switch elevator.Dirn {
	case variabletypes.MD_Up:
		if ordersAbove(elevator, orders) {
			return variabletypes.MD_Up
		} else if ordersBelow(elevator, orders) {
			return variabletypes.MD_Down
		} else {
			return variabletypes.MD_Stop
		}
	case variabletypes.MD_Down:
		if ordersBelow(elevator, orders) {
			return variabletypes.MD_Down
		} else if ordersAbove(elevator, orders) {
			return variabletypes.MD_Up
		} else {
			return variabletypes.MD_Stop
		}
	case variabletypes.MD_Stop:
		if ordersAbove(elevator, orders) {
			return variabletypes.MD_Up
		} else if ordersBelow(elevator, orders) {
			return variabletypes.MD_Down
		} else {
			return variabletypes.MD_Stop
		}
	}
	return variabletypes.MD_Stop
}

func CheckForStop(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix) bool{
	switch elevator.Dirn {
	case variabletypes.MD_Down:
		return (orders[elevator.Floor][1] || orders[elevator.Floor][2] || !ordersBelow(elevator, orders))
	case variabletypes.MD_Up:
		return (orders[elevator.Floor][0] || orders[elevator.Floor][2] || !ordersAbove(elevator, orders))
	case variabletypes.MD_Stop:
		return true
	}
	return false
}


/*
func DelegateOrder(elevMap variabletypes.AllElevatorInfo, buttonEvent variabletypes.ButtonEvent) string {
	AllElevMap := utilities.CreateMapCopy(elevMap)
	currentIP := invalidIP
	currentDuration := 0
	for id, info := range AllElevMap {
		// hvis heisen er i live
		currentElevator = AllElevMap[i]
		elevDuration = costfunction.timeToServeRequest(currentElevator, buttonEvent)
		if elevDuration <= currentDuration {
			currentDuration = elevDuration
			currentIP = i
		}
	}
	return currentIP
}
*/