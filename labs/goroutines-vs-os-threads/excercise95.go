package main

import(
	"os"
	"time"
	"fmt"
	"strconv"
)
var contador int
func main(){
	contador=0
	ida:=make(chan int)
	vuelta:=make(chan int)
	timer := time.NewTimer(1 * time.Second)
	go func(){
		for{
		select {
            case v := <-ida:
                contador++
                fmt.Println(contador)
                vuelta <- v
            }
		}
	}()
	go func(){
		for{
		select {
            case v := <-vuelta:
                contador++
                fmt.Println(contador)
                ida <- v
            }
		}
	}()
    <-timer.C
	ida<-1
	<-vuelta
	f, err := os.Create("Report.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    str:="Messages sent: "+strconv.Itoa(contador)
    fmt.Println(str)
    l, err := f.WriteString(str)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
    }
    fmt.Println(l)
}