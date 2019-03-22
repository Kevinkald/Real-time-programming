package main
import(
	//"fmt"
	"runtime"
	"./variabletypes"
	//"time"
	"./network"
	//"./config"
	"./queuedistribution"
	"./fsm/elevio"
	//"./fsm/fsmdummy"
	"./fsm"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	go queuedistribution.Queuedistribution(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh,buttonsCh,removeOrderCh,ordersCh,elevatorObjectCh)

	go elevio.PollButtons(buttonsCh)

	go fsm.Fsm(ordersCh, elevatorObjectCh, removeOrderCh)

	select{}
}