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
	N_Elevators int = 2

	PeerPort int = 20012
	BroadcastPort int = 30012


	TRAVEL_TIME int = 3 
    DOOR_OPEN_TIME int = 2
)