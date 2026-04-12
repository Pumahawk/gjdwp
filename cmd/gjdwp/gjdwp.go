package main

import (
	"log"
	"net"

	"github.com/Pumahawk/gjdwp/lib/jdwp"
)

func main() {
	add := "localhost:5005"
	tcpConn, err := net.Dial("tcp", add)
	if err != nil {
		log.Fatalf("Unable to connect %s", add)
	}
	conn := jdwp.NewConn(tcpConn)
	done, err := conn.Start()
	if err != nil {
		log.Fatalf("Unable to start jdwp connection: %s", err)
	}
	if r, err := conn.SendVirtualMachineVersion(); err != nil {
		log.Fatalf("Unable to SendVirtualMachineVersion: %s", err)
	} else {
		log.Printf("SendVirtualMachineVersion %v", r)
	}
	<-done
}
