package costfunction

import (
    "../../variabletypes"
    "../../orderlogic"
    "../../config"
    "../utilities"
    "fmt"
)


func TimeToServeRequest (e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent) int { 
    e := CreateMapCopy(e_old)
    e.OrderMatrix[buttonEvnt.Floor][buttonEvnt.Button] = true
    duration := 0
    

    switch e.ElevObj.State {
        case variabletypes.IDLE:
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix);
            if e.ElevObj.Dirn == variabletypes.MD_Stop {
                return duration
            }
            break

        case variabletypes.MOVING:
            duration += config.TRAVEL_TIME/2
            e.ElevObj.Floor += int(e.ElevObj.Dirn)
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)
            break;

        case variabletypes.OPEN:
            duration -= config.DOOR_OPEN_TIME/2
    }

    var arrivedAtRequest bool 
    for {
        if orderlogic.CheckForStop(e.ElevObj, e.OrderMatrix) {
            arrivedAtRequest, e = utilities.Requests_clearAtCurrentFloor(e, buttonEvnt)

            if (arrivedAtRequest) {
                return duration
            }

            duration += config.DOOR_OPEN_TIME
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)
        }


        e.ElevObj.Floor += int(e.ElevObj.Dirn)
        e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)

        duration += config.TRAVEL_TIME
    }
}

func DelegateOrder( elevMap variabletypes.AllElevatorInfo, 
                    PeerInfo variabletypes.PeerUpdate, 
                    buttonEvent variabletypes.ButtonEvent) string {

    if button == variabletypes.BT_Cab {
        return config.ElevatorId 
    }

    peers := PeerInfo.Peers
    button := buttonEvent.Button
    AllElevMap := utilities.CreateMapCopy(elevMap)

    bestId := config.InvalidId
    bestDuration := 99999

    fmt.Println("list of active nodes: ", peers)
    
    for _,id := range peers {
        fmt.Println("entered for loop")
        elevDuration := TimeToServeRequest(AllElevMap[id], buttonEvent)

        if elevDuration <= bestDuration {
            bestDuration = elevDuration
            bestId = id

        }
    }
    return bestId
}
