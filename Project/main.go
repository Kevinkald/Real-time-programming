package main
import(
	//"fmt"
	"runtime"
	//"./config"

	//"./buttons"
	"./network"
)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	//go buttons.Buttons()
	go network.Network()

	


	for {

	}
}