package orderLogic

import (
	. "elevio"
	. "fmt"
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

func chooseNextDirection(singleElevator variableTypes.ElevatorObject, singleElevatorOrders variableTypes.SingleOrderMatrix) {

}