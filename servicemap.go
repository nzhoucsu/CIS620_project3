
package main
 
import (
    "fmt"
    "net"
    "os"
    "strings"
)
 
table_file := "table_file.txt"
 

func main() {
    /* Let's prepare an address at any address at port 10001*/   
    ServerAddr,err := net.ResolveUDPAddr("udp",":23997")
    CheckError(err)
 
    /* Now listen (misleading in Go);  actually bind at selected port */
    ServerCLNS, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerCLNS.Close()
 
    buf := make([]byte, 1024)
 
    for {
        n,addr,err := ServerCLNS.ReadFromUDP(buf)
        fmt.Println("DEBUG--- Received ",string(buf[0:n]), " from ",addr)
        s := strings.Split(string(buf[0:n]), ',')
        if len(s) == 3 // msg from db server
        {
            // Display message from db server (name, IP and port)
            fmt.Println("Received from ", s[1], " : PUT ", s[0], " ", s[2])
            // Write db server name, IP and port into table file
            write_to_table_file(string(buf[0:n]))
            // Reply to db server
            _,err := ServerCLNS.WriteToUDP([]byte("OK"),addr) 
            CheckError(err)
        }
        else if strings.EqualFold('client_request', buf[0:n]) // msg from client
        {
            // return db server IP and port to client.
            data_from_table := read_from_table_file()
            if !strings.EqualFold("no_db_server", data_from_table)
            {
                _,err := ServerCLNS.WriteToUDP(,addr) 
                CheckError(err)
            }            
        }
        else // wrong msg
        {
            fmt.Println("Wrong message from db server or client")
            os.Exit(0)
        }
        
    }
}


/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}


func write_to_table_file(str string){
    defer file.Close()
    // Open or create table file.
    file, err := os.OpenFile(table_file, os.O_CREATE|os.O_WRONLY, 0644)
    CheckError(err error)
    // Write db server name, IP and port into table file.
    _, err = file.WriteString(str)
    CheckError(err)
}


func read_from_table_file() string{
    _, err = os.Stat(path)    
    if os.IsNotExist(err) // create file if not exists
    {
        fmt.Println("service-map server hasn't received any db server submission.")
        return "no_db_server"
    }
    file, err := os.OpenFile(table_file, os.O_RDONLY, 0644)
    CheckError(err)
    data := make([]byte, 100)
    count, err := file.Read(data)
    CheckError(err)
    return data[:count]
}