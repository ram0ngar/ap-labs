// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

//!+
func main() {
	user := flag.String("user", "defaultUser", "User Name String")
	server := flag.String("server", "localhost:9000", "Host address String")
	flag.Parse()
	conn, err := net.Dial("tcp", *server)

	if err != nil {
		log.Fatal(err)
	}
	if blank := strings.TrimSpace(*user) == ""; blank {
		log.Panic("Username cannot be only a space")
	}
	io.WriteString(conn, "<user>"+*user+"\n")
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}