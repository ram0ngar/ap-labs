package main

import(
	"fmt"
	"flag"
	"time"
	"os"
	"strconv"
)

func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func chann(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n
        }
        close(out)
    }()
    return out
}

func main(){
	start:=time.Now()
	param:=flag.Int("channels",5,"number of channels")
	flag.Parse()
	fmt.Println(*param)
	c:=gen(1,2,3)
	for i:=0;i<*param;i++{
		go chann(c)
	}
	elapsed:=time.Since(start)
	f, err := os.Create("Report.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    str:="Channels created: "+strconv.Itoa(*param)+"\nTime to be executed"+elapsed.String()
    l, err := f.WriteString(str)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
    }
    fmt.Println(l)
	fmt.Println(elapsed)
}
