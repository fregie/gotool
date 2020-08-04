package main

import (
	"log"

	"github.com/fregie/gotool/fperf"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fperf.TCPServe(":1201")
}
