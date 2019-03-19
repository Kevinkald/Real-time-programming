package main
import(
	//"fmt"
	"runtime"
	//"./variabletypes"
	//"time"
	//"./buttons"
	//"./network"
	//"./config"
	//"./queuedistribution"
	//"./fsm/elevio"
	//"./fsm/fsmdummy"
	"./fsm/fsm"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Channels between Queuedistributor and Network module
	//peerUpdateCh := make(chan variabletypes.PeerUpdate)
	//networkMessageCh := make(chan variabletypes.AllElevatorInfo)
	//networkMessageBroadcastCh := make(chan variabletypes.AllElevatorInfo)

	//Channel between FSM and Queuedistributor module
	//Insert here
	ordersCh := make(chan variableTypes.SingleOrderMatrix)
	elevatorObjectCh := make(chan variableTypes.ElevatorObject)
	removeOrderCh := make(chan int)

	//Channel between Buttons and Queuedistributor module
	//buttonsCh := make(chan variabletypes.ButtonEvent)
	reachedFloorCh := make(chan int)						// NB! implement!?

	//go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	//go queuedistribution.Queuedistribution(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh,buttonsCh)

	//go elevio.PollButtons(buttonsCh)

	go fsm.Fsm(ordersCh, elevatorObjectCh, removeOrderCh, reachedFloorCh)

	for{}
}