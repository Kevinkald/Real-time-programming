package config
import "os"
import "fmt"
import "time"

var ElevatorPort string
var ElevatorId string

func ConfigInit(){
	ElevatorId = os.Args[1]
	ElevatorPort = os.Args[2]
	fmt.Println(ElevatorPort)
}

const (
	//Interval and timeout in ms
	Interval = 15*time.Millisecond
	Timeout = 1000*time.Millisecond
	InvalidId string = "0"
	NFloors int = 4
	NButtons int = 3
	NElevators int = 3

	PeerPort int = 17563
	BroadcastPort int = 17564

	TRAVEL_TIME int = 3 
    DOOR_OPEN_TIME int = 2
)