package main
import(
	"fmt"
	"runtime"
	"./variabletypes"
	//"time"
	//"./buttons"
	//"./network"
	"./config"
	//"./queuedistribution"
	"./fsm/elevio"
	//"./fsm/fsmdummy"
	"./fsm"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Channels between Queuedistributor and Network module
	//peerUpdateCh := make(chan variabletypes.PeerUpdate)
	//networkMessageCh := make(chan variabletypes.AllElevatorInfo)
	//networkMessageBroadcastCh := make(chan variabletypes.AllElevatorInfo)

	//Channel between FSM and Queuedistributor module
	//Insert here
	ordersCh := make(chan variabletypes.SingleOrderMatrix)
	elevatorObjectCh := make(chan variabletypes.ElevatorObject)
	removeOrderCh := make(chan variabletypes.ElevatorObject)  // variabletypes.ButtonEvent!!!!!!!!!!!!!!!!!!!!!!!

	//Channel between Buttons and Queuedistributor module
	buttonsCh := make(chan variabletypes.ButtonEvent)
	//reachedFloorCh := make(chan int)						// NB! implement!?

	//go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	//go queuedistribution.Queuedistribution(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh,buttonsCh)

	//go elevio.PollButtons(buttonsCh)

	go fsm.Fsm(ordersCh, elevatorObjectCh, removeOrderCh)
	go elevio.PollButtons(buttonsCh)

	var orders variabletypes.SingleOrderMatrix 
	fmt.Println("Here in main")
	for{
		select{
		case msg := <-buttonsCh:
			orders[msg.Floor][msg.Button] = true
			ordersCh <- orders
			fmt.Println("New order")
		case msg := <-removeOrderCh:
			for f := 0; f < config.N_Buttons; f++ {
				orders[msg.Floor][f] = false
			}
			ordersCh <- orders
		}
	}
}