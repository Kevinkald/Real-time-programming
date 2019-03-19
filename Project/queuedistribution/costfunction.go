package costfunction

import (
    "../variabletypes"
)

func timeToServeRequest (e_old variabletypes.SingleElevatorInfo, button variabletypes.ButtonType, floor int) int { 
    var e SingleElevatorInfo = e_old
    e.OrderMatrix[f][b] = 1
    var behavior = // definere denne!!!!

    arrivedAtRequest := 0;

    foo := func ifEqual(inner_b ButtonType, inner_f int) {
        if inner_b == b && inner_f == f {
            arrivedAtRequest = 1;
        }
    }

    duration := 0;

    select {
    case e.behaviour = EB_Idle:
        e.dirn = requests_chooseDirection(e);
        if dirn == D_Stop {
            return duration
        }
        break
    case e.behaviour = EB_Moving:
        duration += TRAVEL_TIME/2
        e.floor += e.dirn
        break;
    case e.behaviour = EB_DoorOpen:
        duration -= DOOR_OPEN_TIME/2

    }

    for {
        if requests_shouldStop(e) {
            e = requests_clearAtCurrentFloor(e, foo(inner_b, inner_f))
            if arrivedAtRequest {
                return duration
            }
            duration += DOOR_OPEN_TIME
            e.dirn = requests_chooseDirection(e)
            }

        e.floor +=e.direction
        duration += TRAVEL_TIME
    }
}















