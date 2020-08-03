package fperf

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
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
	Conn net.Conn
	// CtrlConn     net.Conn
	TestDuration time.Duration
	Stat         *freconn.Stat
	PeerStat     *freconn.Stat
	Role         RoleType
	Finish       bool
}

func NewSender(conn net.Conn, duration time.Duration) *Fperf {
	f := &Fperf{
		Conn:         conn,
		TestDuration: duration,
		Stat:         freconn.NewStat(),
		Role:         Sender,
	}
	return f
}

func NewReceiver(conn net.Conn) *Fperf {
	f := &Fperf{
		Conn: conn,
		Stat: freconn.NewStat(),
		Role: Receiver,
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
	err := f.sendStartOrFin(START)
	if err != nil {
		return err
	}
	// go f.Stat.RunBandwidthIn1()
	// defer f.Stat.StopBandwidthIn1()
	nc := freconn.UpgradeConn(f.Conn)
	nc.EnableStat(f.Stat)
	f.Conn = nc
	KBBuffer := make([]byte, 1024)
	header := &CtrlHeader{
		CtrlType:      DATA,
		PayloadLength: uint32(len(KBBuffer)),
		Data:          KBBuffer,
	}
	dataPacket := header.Pack()
	start := time.Now()
	for time.Now().Before(start.Add(f.TestDuration)) {
		_, err := f.Conn.Write(dataPacket)
		if err != nil {
			log.Printf("Write error: %s", err)
			break
		}
	}
	f.Conn.SetWriteDeadline(time.Time{})
	// log.Printf("send FIN")
	err = f.sendStartOrFin(FIN)
	if err != nil {
		log.Printf("send error: %s", err)
		return err
	}
	// log.Printf("recv Stat")
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
	start, err := f.recvStart()
	if err != nil {
		return err
	}
	// log.Printf("recv START")
	f.Stat.Reset()
	nc := freconn.UpgradeConn(f.Conn)
	nc.EnableStat(f.Stat)
	f.Conn = nc
	for {
		header, err := PacketFromReader(f.Conn)
		if err != nil {
			log.Printf("recv error: %s", err)
			return err
		}
		if header.CtrlType == DATA {
			continue
		}
		if header.CtrlType == FIN {
			// log.Printf("recv FIN")
			ts := int64(binary.BigEndian.Uint64(header.Data))
			f.TestDuration = time.Unix(ts, 0).Sub(start)
			break
		}
	}
	err = f.sendStat()
	if err != nil {
		return err
	}
	f.Finish = true
	return nil
}
