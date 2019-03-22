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

func IfEqual(buttonEvnt variabletypes.ButtonEvent, currentState variabletypes.ElevatorObject) int{
    if buttonEvnt.Floor == currentState.Floor {
        return 1
    }
    return 0
}

func Requests_clearAtCurrentFloor(e_old variabletypes.SingleElevatorInfo, buttonEvnt variabletypes.ButtonEvent, onClearedRequest func(buttonEvnt variabletypes.ButtonEvent, currentState variabletypes.ElevatorObject) int) (int, variabletypes.SingleElevatorInfo) {
    var e variabletypes.SingleElevatorInfo = e_old
    currentFloor := e.ElevObj.Floor

    //var btn variabletypes.ButtonEvent.Button := 0
    onCleared := 0
    for btn := 0; btn < config.K_Buttons; btn ++ {
        if e.OrderMatrix[currentFloor][btn] != 0 {   // if there is an order 
            e.OrderMatrix[currentFloor][btn] = 0;    // clear it
            onCleared = onClearedRequest(buttonEvnt, e.ElevObj)


            /*
            if(onClearedRequest(buttonEvnt, e.ElevObj) != 0) { // Hvordan sjekke om vi skal kjøre denne?
                onCleared = onClearedRequest(buttonEvnt, e.ElevObj)
            } 
            */
        }
    } 
    return onCleared, e;
    
}