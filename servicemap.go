
package main
 
import (
    "fmt"
    "net"
    "os"
    "strings"
)
 

func main() {
    // Set map table
    reg_map := make(map[string] string)
    // Set UDP socket
    ServerAddr,err := net.ResolveUDPAddr("udp",":23997")
    CheckError(err)
    ServerCLNS, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerCLNS.Close()
    // Communicate with db server and client
    buf := make([]byte, 1024)
    for {
        n,addr,err := ServerCLNS.ReadFromUDP(buf)
        CheckError(err)
        s := strings.Split(string(buf[0:n]), ",")
        if len(s) == 3 { // msg from db server        
            // Display message from db server (name, IP and port)
            fmt.Println("Received from", s[1], ": PUT", s[0], s[2])
            // Store db server name, IP and port into table map
            reg_map["hostname"] = s[0]
            reg_map["IP_port"]  = strings.Join([]string{s[1],s[2]}, ":")
            // Reply to db server
            _,err := ServerCLNS.WriteToUDP([]byte("OK"),addr)
            CheckError(err)
        } else{
            addr_s := strings.Split(string(buf[0:n]), ":")
            fmt.Println("Received from", addr_s[0], ": GET", reg_map["hostname"]);
            _,err := ServerCLNS.WriteToUDP([]byte(reg_map["IP_port"]),addr)  
            CheckError(err)
        } 
    }
}


func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

