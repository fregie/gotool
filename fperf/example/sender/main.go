package main

import (
	"flag"
	"log"

	"github.com/fregie/gotool/fperf"
)

var (
	addr = flag.String("a", "127.0.0.1:1030", "addr")
	dura = flag.Int("t", 5, "time")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	result, err := fperf.TCPClientRecv(*addr)
	if err != nil {
		log.Fatal(err)
	}
	result.Print()
}
