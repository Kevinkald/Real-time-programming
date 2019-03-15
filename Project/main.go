package main
import(
	"fmt"
	"runtime"
	"./variabletypes"
	"time"
	//"./buttons"
	"./network"
	"./config"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Channels between Queuedistributor and network module
	peerUpdateCh := make(chan variabletypes.PeerUpdate)
	networkMessageCh := make(chan variabletypes.NetworkMessage)
	networkMessageBroadcastCh := make(chan variabletypes.NetworkMessage)

	go network.Network(peerUpdateCh,networkMessageCh,networkMessageBroadcastCh)

	go func(){
		elevator1 := variabletypes.ElevatorObject{4,-1}
		elevator2 := variabletypes.ElevatorObject{3, 1}
		orders := variabletypes.OrderMatrix{}
		for i:= 0; i<config.M_FLOORS; i++{
			for j:= 0; j<config.N_ELEVATORS; j++{
				orders[i][j] = false
			}
		}

		elevators := variabletypes.Elevatrs{}

		elevators["1"] = {variabletypes.Floor: 4,variabletypes.Dirn: -1}
		elevators["2"] = {3, 1}

		msg := variabletypes.NetworkMessage{Elevators: elevators ,Orders: orders}

		//Keep sending message over to "queuedistributor"
		for {
			networkMessageCh<- msg
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case k := <-peerUpdateCh:
			fmt.Println(k.Peers)

		case m := <-networkMessageCh:
			fmt.Println(m)
		}
	}
}