package freconn

import (
	"log"
	"time"
)

type Stat struct {
	Rx      uint64 `json:"rx"`
	Tx      uint64 `json:"tx"`
	stopCh1 chan bool
	In1     StatStatus `json:"in1"`
	In10    StatStatus `json:"in10"`
}

type StatStatus struct {
	Time        time.Time
	Rx          uint64 `json:"rx"`
	Tx          uint64 `json:"tx"`
	BandwidthRx uint64 `json:"bandwidth_rx"`
	BandwidthTx uint64 `json:"bandwidth_tx"`
}

func (s *StatStatus) Reset() {
	s.Time = time.Now()
	s.Rx = 0
	s.Tx = 0
	s.BandwidthRx = 0
	s.BandwidthTx = 0
}

func NewStat() *Stat {
	s := &Stat{
		Rx: 0,
		Tx: 0,
		In1: StatStatus{
			Time:        time.Now(),
			Rx:          0,
			Tx:          0,
			BandwidthRx: 0,
			BandwidthTx: 0,
		},
		In10: StatStatus{
			Time:        time.Now(),
			Rx:          0,
			Tx:          0,
			BandwidthRx: 0,
			BandwidthTx: 0,
		},
	}

	return s
}

func (s *Stat) RunBandwidthIn1() {
	s.stopCh1 = make(chan bool)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			s.In1.BandwidthRx = s.Rx - s.In1.Rx
			s.In1.BandwidthTx = s.Tx - s.In1.Tx
			s.In1.Rx = s.Rx
			s.In1.Tx = s.Tx
			s.In1.Time = time.Now()

			log.Printf("[1s]RX: %d kbps", s.In1.BandwidthRx/1024)
			log.Printf("[1s]TX: %d kbps", s.In1.BandwidthTx/1024)
		case <-s.stopCh1:
			return
		}
	}
}

func (s *Stat) StopBandwidthIn1() {
	if s.stopCh1 != nil {
		s.stopCh1 <- true
	}
}

func (s *Stat) RunBandwidthIn10() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		s.In10.BandwidthRx = (s.Rx - s.In10.Rx) / 10
		s.In10.BandwidthTx = (s.Tx - s.In10.Tx) / 10
		s.In10.Rx = s.Rx
		s.In10.Tx = s.Tx
		s.In10.Time = time.Now()
		// log.Printf("[10s]RX: %dbps", s.in10.bandwidthRx)
		// log.Printf("[10s]TX: %dbps", s.in10.bandwidthTx)
	}
}

func (s *Stat) AddRx(len uint64) {
	s.Rx += len
}

func (s *Stat) AddTx(len uint64) {
	s.Tx += len
}

func (s *Stat) Bandwidth1() (r, t uint64, lastTime time.Time) {
	r = s.In1.BandwidthRx
	t = s.In1.BandwidthTx
	lastTime = s.In1.Time
	return
}

func (s *Stat) Bandwidth10() (r, t uint64, lastTime time.Time) {
	r = s.In10.BandwidthRx
	t = s.In10.BandwidthTx
	lastTime = s.In10.Time
	return
}

func (s *Stat) Reset() {
	s.Rx = 0
	s.Tx = 0
	s.In1.Reset()
	s.In10.Reset()
}
