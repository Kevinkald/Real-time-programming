package costfunction

import (
    "../../variabletypes"
    //"../../orderlogic"
    "../../config"
    //"../utilities"
    //"math"
    //"strconv"
    //"fmt"
    "math"
)

func DelegateOrder( elevMap variabletypes.AllElevatorInfo, 
                    PeerInfo variabletypes.PeerUpdate, 
                    buttonEvent variabletypes.ButtonEvent) string {
    
    //It its a cab call delegate to this elevator
    if (buttonEvent.Button==variabletypes.BT_Cab){
        return config.ElevatorId
    }

    peers := PeerInfo.Peers
    costs := make(map[string]int)

    for _,Peer := range peers{
        costs[Peer] = CalculateCost(elevMap[Peer], buttonEvent)
    }

    bestoption := "InvalidId"
    bestoptioncost := 9999

    for id, cost := range costs{
        if (cost <= bestoptioncost){
            bestoptioncost = cost
            bestoption = id
        }
    }

    return bestoption
}

func CalculateCost(elevator variabletypes.SingleElevatorInfo, buttonEvent variabletypes.ButtonEvent) int{

    orders := elevator.OrderMatrix
    currentfloor := elevator.ElevObj.Floor
    motordirection := elevator.ElevObj.Dirn
    //state := elevator.ElevObj.State
    orderedfloor := buttonEvent.Floor

    cost := 0

    //Punish distance to order
    difference := int(math.Abs(float64(currentfloor-orderedfloor)))
    cost += difference*10
    //cost += difference * 10

    //Punish # orders
    for floor := 0; floor < config.N_Floors; floor++{

        for button := 0; button < config.N_Buttons; button++{

            if (orders[floor][button]){
                cost += 10
            }
        }
    }

    //Punish direction
    //If the order is in the opposite direction punish!
    /*if (int(currentfloor) + int(motordirection) > orderedfloor){
        cost += 10*difference
    }*/

    if (motordirection==variabletypes.MD_Up){
        if (int(currentfloor)+int(motordirection) > orderedfloor){
            cost += 10/2*difference
        }
    } else if (motordirection==variabletypes.MD_Down){
        if (int(currentfloor)+int(motordirection) < orderedfloor){
            cost += 10/2*difference
        }
    }

 

    return cost
}