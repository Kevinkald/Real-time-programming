package costfunction

import (
    "../../variabletypes"
    "../../orderlogic"
    "../utilities"
    "math"
)


func TimeToServeRequest (e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent) int { 
    button := buttonEvnt.Button
    floor := buttonEvnt.Floor

    var e variabletypes.SingleElevatorInfo = e_old
    e.OrderMatrix[floor][button] = 1

    arrivedAtRequest := 0
    duration := 0
    TRAVEL_TIME := 3 
    DOOR_OPEN_TIME := 2

    switch e.ElevObj.State {

        case variabletypes.IDLE:
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e);
            if e.ElevObj.Dirn == variabletypes.MD_Stop {
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
            temp, _ = utilities.Requests_clearAtCurrentFloor(e, buttonEvnt)
            if temp {
                arrivedAtRequest = 1
            }
            
            _, e = utilities.Requests_clearAtCurrentFloor(e, buttonEvnt)
            if arrivedAtRequest == 1 {
                return duration
            }

            duration += DOOR_OPEN_TIME
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e)
        }

        e.ElevObj.Floor += int(e.ElevObj.Dirn)
        duration += TRAVEL_TIME
    }
}

//func DelegateOrder(elevMap variabletypes.AllElevatorInfo, listOfPeers variabletypes.PeerUpdate.Peers, buttonEvent variabletypes.ButtonEvent, myID variabletypes.NetworkMsg.Id ) string {
func DelegateOrder(elevMap variabletypes.AllElevatorInfo, structOfPeers variabletypes.PeerUpdate, buttonEvent variabletypes.ButtonEvent, myID string ) string {
    structOfPeers.Peers := listOfPeers
    AllElevMap := utilities.CrelistOfPeersateMapCopy(elevMap)
    currentId := config.InvalidId
    currentDuration := Inf(1)
    button := buttonEvent.Button

    if button == BT_Cab {
        return myID 
    }
    for id := range listOfPeers {

        currentElevator = AllElevMap[id]
        elevDuration = timeToServeRequest(currentElevator, buttonEvent)

        if elevDuration <= currentDuration {
            currentDuration = elevDuration
            currentId = id

        }
    }
    return currentId
}
