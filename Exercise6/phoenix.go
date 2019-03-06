package main

import(
	"fmt"
	"os/exec"
	"net"
	"time"
	"encoding/binary"
)

var counter uint64
var buffer = make([]byte, 8)

func startBackup(){
	(exec.Command("gnome-terminal", "-x", "sh", "-c", "go run phoenix.go")).Run()

	fmt.Println("Backup started")
}

func main(){
	primary := false

	addr, _ := net.ResolveUDPAddr("udp","127.0.0.1:5000")
	
	//Secondary listens on addr
	conn, err := net.ListenUDP("udp",addr)
	if (err != nil){
		fmt.Println("error")
	}

	//Sleep to make sure that no numbers are skipped
	//in transition from secondary->primary
	time.Sleep(50 * time.Millisecond)

	//Secondary loop
	for (!primary){
		//Deadline for getting a message set to 2 seconds
		t := time.Now()
		conn.SetReadDeadline(t.Add(time.Second*2))

		//Reads packet into buffer
		_, _, err := conn.ReadFromUDP(buffer)

		//nil means no error, thus we check if there was an error, i.e we
		//received no message
		if (err != nil){
			primary = true
			startBackup()
		} else {
			//Update counter, ugly convertion from []byte to uint64
			counter = binary.BigEndian.Uint64(buffer)

			//Print received number {1,2,3....}
			fmt.Println("Received msg from primary: ",counter)
		}
	} 
	conn.Close()

	//Primary loop
	broadcast_conn, _ := net.DialUDP("udp",nil, addr)
	for (primary){
		counter++

		binary.BigEndian.PutUint64(buffer, counter)

		_, _ = broadcast_conn.Write(buffer)

		fmt.Println("Sent msg to secondary: ",counter)

		//Sleep to avoid DDOS of server..
		time.Sleep(50 * time.Millisecond)
	}
	broadcast_conn.Close()
}