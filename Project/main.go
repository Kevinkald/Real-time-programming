package main
import(
	"fmt"
	"runtime"
	"./variabletypes"
	//"time"
	//"./buttons"
	"./network"
	//"./config"
	"./queuedistribution"
	//"./fsm/elevio"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Channels between Queuedistributor and Network module
	peerUpdateCh := make(chan variabletypes.PeerUpdate)
	networkMessageCh := make(chan variabletypes.AllElevatorInfo)
	networkMessageBroadcastCh := make(chan variabletypes.AllElevatorInfo)

	//Channel between FSM and Queuedistributor module
	//Insert here

	//Channel between Buttons and Queuedistributor module
	//buttonsCh := make(chan variabletypes.ButtonEvent)

	go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	go queuedistribution.Queuedistribution(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	//go elevio.PollButtons(buttonsCh)

	for {
		select {
		case k := <-peerUpdateCh:
			fmt.Println(k.Peers)

		case m := <-networkMessageCh:
			fmt.Println(m)
		}
	}
}