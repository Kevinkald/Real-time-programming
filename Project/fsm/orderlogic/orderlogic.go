package orderlogic

import(
	"../../config"
	"../../variabletypes"
)

func ChooseNextDirection(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)variabletypes.MotorDirection{
	switch elevator.Dirn {
	case variabletypes.MDUp:
		if ordersAbove(elevator, orders) {
			return variabletypes.MDUp
		} else if ordersBelow(elevator, orders) {
			return variabletypes.MDDown
		} else {
			return variabletypes.MDStop
		}
	case variabletypes.MDDown:
		if ordersBelow(elevator, orders) {
			return variabletypes.MDDown
		} else if ordersAbove(elevator, orders) {
			return variabletypes.MDUp
		} else {
			return variabletypes.MDStop
		}
	case variabletypes.MDStop:
		if ordersAbove(elevator, orders) {
			return variabletypes.MDUp
		} else if ordersBelow(elevator, orders) {
			return variabletypes.MDDown
		} else {
			return variabletypes.MDStop
		}
	}
	return variabletypes.MDStop
}

func CheckForStop(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)bool{
	switch elevator.Dirn {
	case variabletypes.MDDown:
		return (orders[elevator.Floor][variabletypes.BTHallDown] || orders[elevator.Floor][variabletypes.BTCab] || !ordersBelow(elevator, orders) || elevator.Floor == 0)
	case variabletypes.MDUp:
		return (orders[elevator.Floor][0] || orders[elevator.Floor][variabletypes.BTCab] || !ordersAbove(elevator, orders) || elevator.Floor == (config.NFloors - 1))
	case variabletypes.MDStop:
		return (orders[elevator.Floor][0] || orders[elevator.Floor][variabletypes.BTHallDown] || orders[elevator.Floor][variabletypes.BTCab])
	}
	return false
}

func ordersAbove(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)bool{
	for floor := elevator.Floor + 1; floor < config.NFloors; floor++ {
		for button := 0; button < config.NButtons; button++ {
			if orders[floor][button] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator variabletypes.ElevatorObject, orders variabletypes.SingleOrderMatrix)bool{
	for floor := 0; floor < elevator.Floor; floor++ {
		for button := 0; button < config.NButtons; button++ {
			if orders[floor][button] {
				return true
			}
		}
	}
	return false
}