package main

import( "net"
        "fmt"
        "bufio"
        "strings"
        "os"
      )

func main() {

  // connect to this socket
  var connections [] net.Conn
  var locations [] string
  i:=0
  for _,param:=range os.Args[1:]{
    i++
    palabras:=strings.Split(param,"=")
    locations=append(locations,palabras[0])
    conn,_:=net.Dial("tcp", palabras[1])
    connections=append(connections,conn)
  }
  fmt.Println(locations)
  for { 
    fmt.Println("\033[0;0H")
    for key,con:=range connections{
        // listen for reply
        message, _ := bufio.NewReader(con).ReadString('\n')
        msg:=strings.Replace(message,"\n","",-1)
        fmt.Printf("%s %s\n",locations[key],msg)
      }
  }
}