package fsm

import (
	"../elevio"
	"fmt"
	"../../config"
	"../../variableTypes"
)

func ordersAbove(elevator variableTypes.ElevatorObject, orders variableTypes.SingleOrderMatrix) {
	for floor := elevator.Floor + 1; floor < config.N_Floors; floor++ {
		for btn := 0; btn < config.N_Buttons; btn++ {
			return true
		}
	}
	return false
}

func ordersBelow(elevator variableTypes.ElevatorObject, orders variableTypes.SingleOrderMatrix) {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < config.N_Buttons; btn++ {
			return true
		}
	}
	return false
}

func ChooseNextDirection(singleElevator variableTypes.ElevatorObject, singleElevatorOrders variableTypes.SingleOrderMatrix) {
	switch singleElevator.Dirn {
	case variableTypes.MD_UP:
		return (!ordersBelow(singleElevator, singleElevatorOrders) || singleELevatorOrders[singleElevator.Floor][]
	case variableTypes.MD_DOWN:

	case variableTypes.MD_STOP:


	}
}

func CheckForStop(singleElevator variableTypes.ElevatorObject, singleElevatorOrders variableTypes.SingleOrderMatrix) {
	
}