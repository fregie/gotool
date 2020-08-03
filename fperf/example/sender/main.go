package main

import (
	"log"
	"time"

	"github.com/fregie/gotool/fperf"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	result, err := fperf.TCPClient("45.77.142.97:1201", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	result.Print()
}
