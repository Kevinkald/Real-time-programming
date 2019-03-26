package costfunction

import (
    "../../variabletypes"
    "../../orderlogic"
    "../../config"
    "../utilities"
    //"math"
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
            //e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)
            break;
        case variabletypes.OPEN:
            duration -= config.DOOR_OPEN_TIME/2
    }
    fmt.Println("cont after switch select state")
    for {
        if orderlogic.CheckForStop(e.ElevObj, e.OrderMatrix) {
            fmt.Println("checkforstop finished")
            arrivedAtRequest, e := utilities.Requests_clearAtCurrentFloor(e, buttonEvnt)
            fmt.Println("arrivedatrequest: ",arrivedAtRequest)
            if (arrivedAtRequest) {
                fmt.Println("arrived at request in floor",e.ElevObj.Floor)
                return duration
            }

            duration += config.DOOR_OPEN_TIME
            e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)
        }
        fmt.Println("floor1: ",e.ElevObj.Floor)

        /*floor := e.ElevObj.Floor
        dir := e.ElevObj.Dirn

        if (floor==0 && dir == -1) {
            return duration
        } else if (floor==config.N_Buttons && dir == 1){
            return duration
        }*/

        e.ElevObj.Floor += int(e.ElevObj.Dirn)

        e.ElevObj.Dirn = orderlogic.ChooseNextDirection(e.ElevObj, e.OrderMatrix)

        duration += config.TRAVEL_TIME
        fmt.Println("floor2: ",e.ElevObj.Floor)
    }
}

//func DelegateOrder(elevMap variabletypes.AllElevatorInfo, listOfPeers variabletypes.PeerUpdate.Peers, buttonEvent variabletypes.ButtonEvent, myID variabletypes.NetworkMsg.Id ) string {
func DelegateOrder(elevMap variabletypes.AllElevatorInfo, structOfPeers variabletypes.PeerUpdate, buttonEvent variabletypes.ButtonEvent, myID string ) string {
    listOfPeers := structOfPeers.Peers
    AllElevMap := utilities.CreateMapCopy(elevMap)
    currentId := config.InvalidId

    //currentDuration := int(math.Inf(1))
    currentDuration := 99999
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