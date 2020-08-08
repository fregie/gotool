package main

import (
	"flag"
	"log"
	"time"

	"github.com/fregie/gotool/fperf"
)

var (
	addr = flag.String("a", "127.0.0.1:1030", "addr")
	dura = flag.Int("t", 5, "time")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	start := time.Now()
	result, err := fperf.TCPClientRecv(*addr, time.Duration(*dura)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("test duration: %d s", int(time.Since(start).Seconds()))
	result.Print()
}
