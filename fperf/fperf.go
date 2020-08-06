package fperf

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/fregie/gotool/freconn"
)

const (
	uint64Length = 8
)

type RoleType int

const (
	RoleType_Unknown RoleType = iota
	Sender
	Receiver
)

type Fperf struct {
	Conn         net.Conn
	CtrlConn     net.Conn
	TestDuration time.Duration
	Stat         *freconn.Stat
	PeerStat     *freconn.Stat
	Role         RoleType
	Finish       bool
}

func NewSender(conn, ctrlConn net.Conn, duration time.Duration) *Fperf {
	f := &Fperf{
		Conn:         conn,
		CtrlConn:     ctrlConn,
		TestDuration: duration,
		Stat:         freconn.NewStat(),
		Role:         Sender,
	}
	return f
}

func NewReceiver(conn, ctrlConn net.Conn) *Fperf {
	f := &Fperf{
		Conn:     conn,
		CtrlConn: ctrlConn,
		Stat:     freconn.NewStat(),
		Role:     Receiver,
	}
	return f
}

func (f *Fperf) GetStat() *freconn.Stat {
	return f.Stat
}

func (f *Fperf) GetPeerStat() *freconn.Stat {
	return f.PeerStat
}

func (f *Fperf) RunSender() error {
	if f.Role != Sender {
		return errors.New("not a sender")
	}
	// log.Printf("send START")
	err := f.sendStart()
	if err != nil {
		return err
	}
	// go f.Stat.RunBandwidthIn1()
	// defer f.Stat.StopBandwidthIn1()
	nc := freconn.UpgradeConn(f.Conn)
	nc.EnableStat(f.Stat)
	f.Conn = nc
	KBBuffer := make([]byte, 1024)
	start := time.Now()
	f.Conn.SetWriteDeadline(start.Add(f.TestDuration))
	for {
		_, err := f.Conn.Write(KBBuffer)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				break
			}
			log.Printf("Write error: %s", err)
			return err
		}
	}
	// log.Printf("send FIN")
	err = f.sendFIN()
	if err != nil {
		log.Printf("send error: %s", err)
		return err
	}
	// log.Printf("recv Stat")
	err = f.sendStat()
	if err != nil {
		return err
	}
	err = f.recvStat()
	if err != nil {
		log.Printf("send error: %s", err)
		return err
	}
	f.Finish = true
	return nil
}

func (f *Fperf) RunReceiver() error {
	if f.Role != Receiver {
		return errors.New("not a receiver")
	}
	data, err := f.recvStart()
	if err != nil {
		return err
	}
	start := time.Unix(data.Start, 0)
	f.TestDuration = time.Duration(data.TestDuration) * time.Second
	// log.Printf("recv START")
	f.Stat.Reset()
	nc := freconn.UpgradeConn(f.Conn)
	nc.EnableStat(f.Stat)
	f.Conn = nc
	f.Conn.SetReadDeadline(start.Add(f.TestDuration))
	buffer := make([]byte, 65535)
	for {
		_, err := f.Conn.Read(buffer)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				break
			}
			log.Printf("recv error: %s", err)
			return err
		}
	}
	err = f.recvFIN()
	if err != nil {
		return err
	}
	err = f.sendStat()
	if err != nil {
		return err
	}
	err = f.recvStat()
	if err != nil {
		log.Printf("send error: %s", err)
		return err
	}
	f.Finish = true
	return nil
}
