package config
import "os"

var Port string
var ID string


func Init(){
	Port = os.Args[1]
	ID = os.Args[2]
}

// ID ikke ferdig implementert
const (

	ElevatorId string = "1"
	InvalidId string = "0"

	N_Floors int = 4
	N_Buttons int = 3
	N_Elevators int = 2

	PeerPort int = 20012
	BroadcastPort int = 30012
	SimulatorPort string = "localhost:16001"
	HardwarePort string = "localhost:15657"
)