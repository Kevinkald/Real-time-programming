package config
import "os"
import "fmt"

var ElevatorPort string
var ElevatorId string


func ConfigInit(){
	ElevatorId = os.Args[1]
	ElevatorPort = os.Args[2]
	fmt.Println(ElevatorPort)
}

const (
	InvalidId string = "0"
	N_Floors int = 4
	N_Buttons int = 3
	N_Elevators int = 3

	PeerPort int = 17563
	BroadcastPort int = 17564


	TRAVEL_TIME int = 3 
    DOOR_OPEN_TIME int = 2
)