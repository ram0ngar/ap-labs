// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	usersChan = make(map[string] net.Conn)
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	who := ""

	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		txt:=input.Text()
		if match,_:=regexp.MatchString("<user>.+",txt); match{
			who=strings.Split(txt,">")[1]
			ch <- "irc-server > The user [" + who + "] successfully logged"
			messages <- "irc-server > [" + who + "] has arrived"
			entering<-ch
			usersChan[who]=conn
		}else if txt=="/users"{
			fmt.Fprintf(conn,"irc-server >")
			for key:=range(usersChan){
				fmt.Fprintf(conn,"%s",key)
			}
			fmt.Fprintf(conn,"\n")
		}else if match, _ := regexp.MatchString("^/msg .+ .+", txt); match{
			strimString := strings.Split(txt, " ")
			lenSlice := len(strimString)
			if user, check := usersChan[strimString[1]]; check {
				fmt.Fprintf(user, "%s [privateMSG] > ", who)
				for i := 2; i < lenSlice; i++ {
					fmt.Fprintf(user, " %s ", strimString[i])
				}
				fmt.Fprintf(user, "\n")
			} else {
				fmt.Fprintf(conn, "irc-server > No user (%s) found, use /users to see connected users. \n", strimString[1])
			}
		}else if txt == "/time"{
			loc, _ := time.LoadLocation("America/Mexico_city")
	    	tme := time.Now().In(loc)
			timeFormat := tme.Format("15:04") 
	    	fmt.Fprintf(conn, "irc-server > Local Time: %s %s\n", loc, timeFormat)
	    }else if match, _ := regexp.MatchString("^/user .+$", txt); match{
	    	strimString := strings.Split(txt, " ")
			if user, check := usersChan[strimString[1]]; check {
				fmt.Fprintf(conn, "irc-server > username: %s, IP: %s \n", strimString[1], user.RemoteAddr().String())
			} else {
				fmt.Fprintf(conn, "irc-server > No user (%s) found, use /users to see connected users. \n", strimString[1])
			}
		}else{
			messages <- who + "> " + input.Text()
		}

	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- "irc-server > [" + who + "] has left"
	fmt.Printf("irc-server > [%s] left \n", who)
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	host := flag.String("host", "localhost", "host string")
	port := flag.String("port", "9000", "port string")
	flag.Parse()
	fmt.Printf("irc-server > Simple IRC Server start at %s:%s \n", *host, *port)
	fmt.Printf("irc-server > Ready for receiving new clients \n")
	listener, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
//!-main