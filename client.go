package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"bufio"
	"encoding/binary"
	"unsafe"
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
	// Send request to service-map server to ask db server IP and port
	hostname, err := os.Hostname()
	CheckError(err)
	hostaddrs, err := net.LookupHost(hostname)
	CheckError(err)
	_, err = CLNS.WriteToUDP([]byte(hostaddrs[0]), ServerAddr)
	CheckError(err)
	// Get db server IP and port from service-map server
	buff := make([]byte, 1024)
	n, _, err := CLNS.ReadFromUDP(buff)
	CheckError(err)
	s = strings.Split(string(buff[0:n]), ":")
	fmt.Println("Service provided by", s[0], "at port", s[1])
	// Close UDP socket for communication between service-map server and client
	CLNS.Close()	

	// Set up TCP client socket for connection between db server and client
	conn, err := net.Dial("tcp", string(buff[0:n]))
	CheckError(err)
	fmt.Println("messge out ......") 
	fmt.Println([]byte(hostaddrs[0])) 
	n, err = conn.Write([]byte(hostaddrs[0]))
	fmt.Println("send out",n,"characters")
	n, err = conn.Read(buff)
	fmt.Println("messge back ......") 
	fmt.Println("send back",n,"characters")
	s = strings.Split(string(buff[0:n]), ":")
	fmt.Println(s[0])
	conn.Close()
	
}


func CheckError(err error){
	if err != nil{
		fmt.Println("Error: ", err)
	}
}

func db_oprt(con net.Conn){
	var s string
	var iprt *uint32
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s = scanner.Text()
		iprt = (*uint32)(unsafe.Pointer(&s))
		buf := make([]byte, 1024)
		binary.BigEndian.PutUint32(buf, *iprt)
		con.Write(buf)
		// n, err = conn.Read(buff)
		// fmt.Println(string(buff[0:n]))
	}
}
