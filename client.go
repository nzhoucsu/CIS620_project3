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
	hostname, err := os.Hostname()
	CheckError(err)
	hostaddrs, err := net.LookupHost(hostname)
	CheckError(err)
	_, err = CLNS.WriteToUDP([]byte(hostaddrs[0]), ServerAddr)
	CheckError(err)

	// Get db server IP and port
	buff := make([]byte, 1024)
	n, _, err := CLNS.ReadFromUDP(buff)
	s = strings.Split(string(buff[0:n]), ":")
	fmt.Println("Service provided by", s[0], "at port", s[1])
}


func CheckError(err error){
	if err != nil{
		fmt.Println("Error: ", err)
	}
}
