package main

import (
	"log"
	"time"

	"github.com/fregie/gotool/fperf"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fperf.TCPSendServe(":1201", 5*time.Second)
}
