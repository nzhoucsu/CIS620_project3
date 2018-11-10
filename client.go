package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
	// Get service-map server IP
	addrs, err := net.LookupHost(os.Args[1])
	CheckError(err)
	fmt.Println("DEBUG---")
	fmt.Println("service-map server name is", os.Args[1])
	fmt.Println("service-map server IP is", addrs[0])

	// Set up UDP socket to communicate with service-map server
	s := []string{addrs[0], "23997"}
	ServerAddr, err := net.ResolveUDPAddr("udp", strings.Join(s, ":"))
	CheckError(err)
	LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
	CheckError(err)
	CLNS, err := net.ListenUDP("udp", LocalAddr)
	CheckError(err)
	defer CLNS.Close()

	// Send request to service-map server to ask db server IP and port
	_, err = CLNS.WriteToUDP([]byte("client_request"), ServerAddr)
	CheckError(err)

	// Get db server IP and port
	buff := make([]byte, 1024)
	n, _, err := CLNS.ReadFromUDP(buff)
	fmt.Println("Get from service-map server:", buff[0:n])
}


func CheckError(err error){
	if err != nil{
		fmt.Println("Error: ", err)
	}
}