package utilities

import(
	"../../variabletypes"
	"../../config"
	"strconv"
	"fmt"
	//"time"
)

func CreateMapCopy(elevMap variabletypes.AllElevatorInfo) variabletypes.AllElevatorInfo {
	copyMap := make(variabletypes.AllElevatorInfo)
	for key, value := range elevMap {
		copyMap[key] = value
	}
	return copyMap
}

func InitMap()variabletypes.AllElevatorInfo{
	elevMap := make(map[string]variabletypes.SingleElevatorInfo)
	for id := 1; id <= config.N_Elevators; id++{
		id_string := strconv.Itoa(id)
		elevMap[id_string] = variabletypes.SingleElevatorInfo{}
	}
	return elevMap
}

func PrintMap(a variabletypes.AllElevatorInfo){
		for id := 1; id <= config.N_Elevators; id++{
			id_string := strconv.Itoa(id)
			fmt.Println("Elevator id: ",id_string)
			for floor := 0; floor < config.N_Floors; floor++{
				fmt.Println(a[id_string].OrderMatrix[floor])
			}
			fmt.Println("State", a[id_string].ElevObj.State)
			fmt.Println("Floor", a[id_string].ElevObj.Floor)
			fmt.Println("Dirn", a[id_string].ElevObj.Dirn)
		}
		//time.Sleep(200*time.Millisecond)
}

/*
func IfEqual(requestedButton variabletypes.ButtonType, requestedFloor variabletypes.ButtonType, currentState variabletypes.ElevatorObject) bool {
	currentState.Floor := inner_f
	// currentState.Button := inner_b

	if inner_b == requestedButton && inner_f == requestedFloor {
        return true
    }
    return false
}

func Requests_clearAtCurrentFloor(e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent, onClearedRequest func(requestedButton variabletypes.ButtonType, requestedFloor variabletypes.ButtonType, currentState variabletypes.ElevatorObject) int) (int, variabletypes.SingleElevatorInfo) {
    var e variabletypes.SingleElevatorInfo = e_old
    currentFloor := e.ElevObj.Floor

    var btn variabletypes.ButtonEvent := 0
    onCleared := false

    for btn = 0; btn < config.N_Buttons; btn ++ {
        if e.OrderMatrix[currentFloor][btn] = true {   // if there is an order 
            e.OrderMatrix[currentFloor][btn] = false;    // clear it
            onCleared = onClearedRequest(btn, buttonEvnt.Floor, e.ElevObj))
        }
    } 
    return onCleared, e;
}

*/

func Requests_clearAtCurrentFloor(e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent) (bool, variabletypes.SingleElevatorInfo) {
    var e variabletypes.SingleElevatorInfo = e_old
    currentFloor := e.ElevObj.Floor
    onCleared := false

    //for btn := 0; btn < config.N_Buttons; btn ++ {
    for btn := variabletypes.BT_HallUp; btn <= variabletypes.BT_Cab; btn ++ {
        if e.OrderMatrix[currentFloor][btn] == true {   // if there is an order 
            e.OrderMatrix[currentFloor][btn] = false;    // clear it
            fmt.Println("oncleared: ", onCleared)
            if (buttonEvnt.Floor == currentFloor)/* && (buttonEvnt.Button == btn)*/ {
            	onCleared = true
            }
            fmt.Println("oncleared after if: ", onCleared)
        }
    } 
    return onCleared, e;
}


func ClearOrder(r int, tmp variabletypes.SingleElevatorInfo) {
    for button := 0; button < config.N_Buttons; button++ {
            tmp.OrderMatrix[r][button] = false      
    }
}









