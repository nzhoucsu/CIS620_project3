package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"bufio"
	"encoding/binary"
	"unsafe"
	"strconv"
	"bytes"
)


func main() {
	// Check input
	if len(os.Args) != 2 {
		fmt.Println("Please input service-map server name.")
		return
	}
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

	// Implement communication with database file on db server.
	db_oprt(conn)
	conn.Close()	
}


func CheckError(err error){
	if err != nil{
		fmt.Println("Error: ", err)
	}
}

func db_oprt(con net.Conn){
	var cmd  int
	var acct int
	var amnt float32
	var iprt *uint32
	var err error
	buf_cmd  := make([]byte, 4)
	buf_acct := make([]byte, 4)
	buf_amnt := make([]byte, 4)
	buf := make([]byte, 30)
	// Preparation for user command input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.Compare(scanner.Text(), "quit") == 0 {
			con.Write([]byte("quit"))
			return
		}
		s := strings.Split(scanner.Text(), " ")
		if strings.Compare(s[0], "query") == 0 { // User requests query 
			cmd = 1000
			iprt = (*uint32)(unsafe.Pointer(&cmd))
			binary.BigEndian.PutUint32(buf_cmd, *iprt)
			acct, err = strconv.Atoi(s[1])
			CheckError(err)
			iprt = (*uint32)(unsafe.Pointer(&acct))
			binary.BigEndian.PutUint32(buf_acct, *iprt)
			s := [][]byte{buf_cmd, buf_acct}
			buf_query := bytes.Join(s, []byte(""))
			con.Write(buf_query)			
		} else if strings.Compare(s[0], "update") == 0 { // User requests update
			cmd = 1001
			iprt = (*uint32)(unsafe.Pointer(&cmd))
			binary.BigEndian.PutUint32(buf_cmd, *iprt)
			acct, err = strconv.Atoi(s[1])
			CheckError(err)
			iprt = (*uint32)(unsafe.Pointer(&acct))
			binary.BigEndian.PutUint32(buf_acct, *iprt)
			value, err := strconv.ParseFloat(s[2], 32)
			CheckError(err)
			amnt = float32(value)
			iprt = (*uint32)(unsafe.Pointer(&amnt))
			binary.BigEndian.PutUint32(buf_amnt, *iprt)
			s := [][]byte{buf_cmd, buf_acct, buf_amnt}
			buf_update := bytes.Join(s, []byte(""))
			con.Write(buf_update)
		} else {
			fmt.Println("Wrong command")
			continue
		}
		n, err := con.Read(buf)
		CheckError(err)
		fmt.Println(string(buf[0:n]))
	}
}
