// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"fmt"
	"log"
	"os"
	"flag"
	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

//!-sema

//!+
func main() {
	depth:=flag.Int("depth",5,"Echale el numerazo")
	flag.Parse()
	fmt.Println(*depth)
	worklist := make(chan []string)
	var n int // number of pending sends to worklist
	var count int

	// Start with the command-line arguments.
	n++
	fmt.Println(os.Args[2:])
	go func() { worklist <- os.Args[2:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		count=0
		for _, link := range list {
			if count>*depth{
				break
			}else{
				if !seen[link] {
					seen[link] = true
					n++
					go func(link string) {
						worklist <- crawl(link)
					}(link)
				}
			}
			count++
		}
	}
}

//!-