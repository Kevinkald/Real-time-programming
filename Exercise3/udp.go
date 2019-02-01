package main

import(
	"fmt"
	"runtime"
	"net"
	"log"
	"time"
)

func receive(portnumber int) {
	addr := net.UDPAddr{Port: portnumber}
	pc, err := net.ListenUDP("udp", &addr)
	if err != nil{
		fmt.Println(err)
	}
	defer pc.Close()
	buffer := make([]byte, 1024)
	for {
		n, err2 := pc.Read(buffer)
		if err2!= nil{
			fmt.Println(err)
		}
		fmt.Println(string(buffer[:n]))
	}
}


func sender(finished chan<- bool){
	//Conn, err := net.Dial("udp", "10.100.23.242:20012")
	Conn, err := net.Dial("udp","255.255.255.255:20012")
	if err != nil {
			log.Fatal(err)
		}
	defer Conn.Close()
	for {
		Conn.Write([]byte("Hello!!!!!!!!!!"))
		time.Sleep(500 * time.Millisecond)

	}

	finished <- true
}


/*
	IP ADDRESS SERVER: 10.100.23.242:58938
*/
func main(){

	runtime.GOMAXPROCS(runtime.NumCPU())
	
	finished := make(chan bool)

	go sender(finished)
	go receive(20012)

	<- finished
	
}

