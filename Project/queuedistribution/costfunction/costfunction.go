package costfunction

import (
    "../../variabletypes"
    "../../orderlogic"
    "../../config"
    "../utilities"
    "math"
    //"strconv"
    "fmt"
)


func TimeToServeRequest (e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent) int { 
    var e variabletypes.SingleElevatorInfo = e_old
    e.OrderMatrix[buttonEvnt.Floor][buttonEvnt.Button] = true
    duration := 0
    fmt.Println(e.ElevObj.State)
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
            break;
        case variabletypes.OPEN:
            duration -= config.DOOR_OPEN_TIME/2
    }

    for {
        if orderlogic.CheckForStop(e.ElevObj, e.OrderMatrix) {
            arrivedAtRequest, e := utilities.Requests_clearAtCurrentFloor(e, buttonEvnt)

            if arrivedAtRequest {
                return duration
            }

            duration += config.DOOR_OPEN_TIME
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)
        }

        e.ElevObj.Floor += int(e.ElevObj.Dirn)
        duration += config.TRAVEL_TIME
    }
}

//func DelegateOrder(elevMap variabletypes.AllElevatorInfo, listOfPeers variabletypes.PeerUpdate.Peers, buttonEvent variabletypes.ButtonEvent, myID variabletypes.NetworkMsg.Id ) string {
func DelegateOrder(elevMap variabletypes.AllElevatorInfo, structOfPeers variabletypes.PeerUpdate, buttonEvent variabletypes.ButtonEvent, myID string ) string {
    listOfPeers := structOfPeers.Peers
    AllElevMap := utilities.CreateMapCopy(elevMap)
    currentId := config.InvalidId

    currentDuration := int(math.Inf(1))
    button := buttonEvent.Button

    fmt.Println("list of active nodes: ", listOfPeers)

    if button == variabletypes.BT_Cab {
        return myID 
    }
    for _,id := range listOfPeers {

        fmt.Println("entere for loop")
        
        currentElevator := AllElevMap[id]
        elevDuration := TimeToServeRequest(currentElevator, buttonEvent)
        fmt.Println("calculated elevDur to be: ", elevDuration)

        if elevDuration <= currentDuration {
            currentDuration = elevDuration
            currentId = id

        }
    }
    return currentId
}




