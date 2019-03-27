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
	peerUpdateCh := make(chan variabletypes.PeerUpdate,10)
	networkMessageCh := make(chan variabletypes.NetworkMsg,10)
	networkMessageBroadcastCh := make(chan variabletypes.NetworkMsg,10)

	//Channel between FSM and Queuedistributor module
	ordersCh := make(chan variabletypes.SingleOrderMatrix,10)
	elevatorObjectCh := make(chan variabletypes.ElevatorObject,10)
	removeOrderCh := make(chan int,10)


	//Channel between Buttons and Queuedistributor module
	buttonsCh := make(chan variabletypes.ButtonEvent,10)

	//Ch buttons
	elevatorsCh := make(chan variabletypes.AllElevatorInfo,10)
	alivePeersCh := make(chan variabletypes.PeerUpdate,10)

	go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	go queuedistribution.Queuedistribution(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh,buttonsCh,removeOrderCh,ordersCh,elevatorObjectCh,elevatorsCh,alivePeersCh)

	go synchlogic.SynchronizeButtonLamps(elevatorsCh,alivePeersCh)

	go elevio.PollButtons(buttonsCh)

	go fsm.Fsm(ordersCh, elevatorObjectCh, removeOrderCh)

	select{}
}