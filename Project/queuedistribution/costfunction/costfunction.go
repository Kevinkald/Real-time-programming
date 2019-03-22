package costfunction

import (
    "../../variabletypes"
    "../../orderlogic"
    "../utilities"
)


func TimeToServeRequest (e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent) int { 
    button := buttonEvnt.Button
    floor := buttonEvnt.Floor

    var e variabletypes.SingleElevatorInfo = e_old
    e.OrderMatrix[floor][button] = 1
    behaviour := e.ElevObj.State

    arrivedAtRequest := 0
    duration := 0
    TRAVEL_TIME := 10 

    select {
        //case behaviour = variabletypes.IDLE:
        case variabletypes.IDLE:
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e);
            if e.ElevObj.Dirn == MD_Stop {
                return duration
            }
            break
        case variabletypes.MOVING:
            duration += TRAVEL_TIME/2
            e.ElevObj.Floor += e.ElevObj.Dirn
            break;
        case variabletypes.OPEN:
            duration -= DOOR_OPEN_TIME/2

    }

    for {
        if orderlogic.CheckForStop(e) {

            arrivedAtRequest = utilities.IfEqual(floor, e.ElevObj.Floor)
            e = utilities.Requests_clearAtCurrentFloor(e, IfEqual(floor, e.ElevObj.Floor))
            if arrivedAtRequest {
                return duration
            }

            duration += DOOR_OPEN_TIME
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e)
        }

        e.ElevObj.Floor += e.ElevObj.Dirn
        duration += TRAVEL_TIME
    }
}