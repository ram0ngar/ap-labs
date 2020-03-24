/*
Resources:
https://golang.org/pkg/os/



*/

package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
)

var directories,devices,sockets,links,other int = 0,0,0,0,0;

func check(path string, infoX os.FileInfo, errX error)error{

  fileType:= string(os.FileMode.String(infoX.Mode())[0])

  switch fileType{
  case "d":
    directories++
  case "L":
    links++
  case "S":
    sockets++
  case "D":
    devices++
  default:
    other++
  }
  return nil
}

func scanDir(dir string) error {

  err := filepath.Walk(dir,check);
	if err!=nil{
		log.Fatal(err)
	}
  return err
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dir-scan <directory>")
		os.Exit(1)
	}

  scanDir(os.Args[1])

	fmt.Println("Directory Scanner Tool")
	fmt.Println("+---------------+-------+")
	fmt.Println("|Path\t\t|",os.Args[1],"\t|")
	fmt.Println("+---------------+-------+")
	fmt.Println("|Directories\t|",directories,"\t|")
	fmt.Println("|Symbolic links\t|",links,"\t|")
	fmt.Println("|Devices\t|",devices,"\t|")
	fmt.Println("|Sockets\t|",sockets,"\t|")
	fmt.Println("|Other files\t|",other,"\t|")
	fmt.Println("+---------------+-------+")
}