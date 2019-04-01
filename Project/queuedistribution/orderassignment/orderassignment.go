package orderassignment
import (
    "../../variabletypes"
    "../../config"
    "../utilities"
    "math"
)

func DelegateOrder( elevMap variabletypes.AllElevatorInfo, 
                    peerInfo variabletypes.PeerUpdate, 
                    buttonEvent variabletypes.ButtonEvent) string {
    
    if (buttonEvent.Button==variabletypes.BTCab){
        return config.ElevatorId
    }

    peers := peerInfo.Peers
    costs := make(map[string]int)

    for _,peer := range peers{
        costs[peer] = calculateCost(elevMap[peer], buttonEvent)
    }

    optimalElevator := ""
    optimalElevatorCost := 9999

    for elevId, cost := range costs{
        if (cost <= optimalElevatorCost){
            optimalElevatorCost = cost
            optimalElevator = elevId
        }
    }

    return optimalElevator
}

func RedistributeOrders( peers variabletypes.PeerUpdate,
                         elevatorMap variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo {
    
    redistributedMap := utilities.CreateMapCopy(elevatorMap)
    var redistributedOrder variabletypes.ButtonEvent

    for _,lostElevatorId := range peers.Lost {
        for floor := 0; floor < config.NFloors; floor++ {
            redistributedOrder.Floor = floor
            for button := variabletypes.BTHallUp; button <= variabletypes.BTHallDown; button++{
                redistributedOrder.Button = button
                
                if (elevatorMap[lostElevatorId].OrderMatrix[floor][button]){
                    newId := DelegateOrder(elevatorMap, peers, redistributedOrder)
                    redistributedMap[newId] = 
                    utilities.SetSingleElevatorMatrixValue(redistributedMap[newId], floor, int(button), true);
                }
            }
        }
    }
    return redistributedMap
}

func calculateCost( elevator variabletypes.SingleElevatorInfo,
                    buttonEvent variabletypes.ButtonEvent) int{

    orders := elevator.OrderMatrix
    currentFloor := elevator.ElevObj.Floor
    motorDirection := elevator.ElevObj.Dirn
    orderedFloor := buttonEvent.Floor

    cost := 0

    //Punish distance to order
    floorDifference := int(math.Abs(float64(currentFloor-orderedFloor)))
    cost += floorDifference*10

    //Punish # orders
    for floor := 0; floor < config.NFloors; floor++{
        for button := 0; button < config.NButtons; button++{
            if (orders[floor][button]){
                cost += 10
            }
        }
    }

    //Punish wrong motor direction
    if (motorDirection==variabletypes.MDUp){
        if (int(currentFloor)+int(motorDirection) > orderedFloor){
            cost += 10/2*floorDifference
        }
    } else if (motorDirection==variabletypes.MDDown){
        if (int(currentFloor)+int(motorDirection) < orderedFloor){
            cost += 10/2*floorDifference
        }
    }

    return cost
}