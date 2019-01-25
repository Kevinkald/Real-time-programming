package main

import(
	. "fmt"
	"runtime"
	"net"
	"log"
	"time"
)

func receive(portnumber string) {
	pc, err := net.ListenPacket("udp",portnumber)
	for {
		
		if err != nil{
			log.Fatal(err)
		}
		defer pc.Close()
	
		buffer := make([]byte, 1024)
		Print(pc.ReadFrom(buffer))

		time.Sleep(500* time.Millisecond)
	}
}


func sender(finished chan<- bool){
	// broadcastIP = #.#.#.255. First three bytes are from the local IP, or just use 255.255.255.255
	//addr := UDPAddr{IP: 10.100.23.242, Port: 20012, Zone: ""}
	//addr = new InternetAddress(broadcastIP, port)
	Conn, err := net.Dial("udp", "10.100.23.242:20012")
	for {
		if err != nil {
			log.Fatal(err)
		}
		defer Conn.Close()
		Conn.Write([]byte("Hello"))
		time.Sleep(500 * time.Millisecond)

	}

	finished <- true
	//sendSock = new Socket(udp) // UDP, aka SOCK_DGRAM
	//sendSock.setOption(broadcast, true)
	//sendSock.sendTo(message, 10.100.23.242:58938)
}


/*
	IP ADDRESS SERVER: 10.100.23.242:58938
*/
func main(){

	runtime.GOMAXPROCS(runtime.NumCPU())
	
	finished := make(chan bool)
	//1
	//receive(":30000")

	//2
	go sender(finished)
	go receive(":20012")

	<- finished
	//go receive()

	//receive()
	
}

