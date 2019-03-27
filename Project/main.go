package main
import(
	//"fmt"
	"runtime"
	"./variabletypes"
	//"time"
	"./network"
	"./config"
	"./queuedistribution"
	"./fsm/elevio"
	//"./fsm/fsmdummy"
	"./fsm"
	"./queuedistribution/synchlogic"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	config.ConfigInit()

	elevio.Init("localhost:" + config.ElevatorPort)

	//Channels between Queuedistributor and Network module
	peerUpdateCh := make(chan variabletypes.PeerUpdate)
	networkMessageCh := make(chan variabletypes.NetworkMsg)
	networkMessageBroadcastCh := make(chan variabletypes.NetworkMsg)

	//Channel between FSM and Queuedistributor module
	ordersCh := make(chan variabletypes.SingleOrderMatrix)
	elevatorObjectCh := make(chan variabletypes.ElevatorObject)
	removeOrderCh := make(chan int)


	//Channel between Buttons and Queuedistributor module
	buttonsCh := make(chan variabletypes.ButtonEvent)

	//Ch buttons
	elevatorsCh := make(chan variabletypes.AllElevatorInfo)
	alivePeersCh := make(chan variabletypes.PeerUpdate)

	go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	go queuedistribution.Queuedistribution(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh,buttonsCh,removeOrderCh,ordersCh,elevatorObjectCh,elevatorsCh,alivePeersCh)

	go synchlogic.SynchronizeButtonLamps(elevatorsCh,alivePeersCh)

	go elevio.PollButtons(buttonsCh)

	go fsm.Fsm(ordersCh, elevatorObjectCh, removeOrderCh)

	select{}
}