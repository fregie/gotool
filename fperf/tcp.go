package fperf

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/brucespang/go-tcpinfo"
)

const RWBufferSize = 1048576

type ConnType uint8

const (
	UnsetConn ConnType = iota
	TransConn
	CtrlConn
)
const FirstDataLen = 5

type TcpPerfServer struct {
	Addr string
}

var connMap sync.Map

func TCPSendServeWithContext(ctx context.Context, Addr string, duration time.Duration) error {
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s", lis.Addr().String())
	connCh := make(chan net.Conn)
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Printf("accept err: %s", err)
				log.Printf("Stop fperf listening")
				return
			}
			conn.(*net.TCPConn).SetWriteBuffer(RWBufferSize)
			conn.(*net.TCPConn).SetDeadline(time.Now().Add(duration + 5*time.Second))
			connCh <- conn
		}
	}()
	for {
		select {
		case conn := <-connCh:
			handleConn(conn, duration)
		case <-ctx.Done():
			lis.Close()
			return nil
		}
	}
}

var TCPListener net.Listener

func TCPSendServe(Addr string, duration time.Duration) error {
	return TCPSendServeWithContext(context.Background(), Addr, duration)
}

func StopTCPServe() error {
	if TCPListener == nil {
		return errors.New("Not running")
	}
	defer func() { TCPListener = nil }()
	return TCPListener.Close()
}

func handleConn(conn net.Conn, duration time.Duration) error {
	data := make([]byte, FirstDataLen)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		conn.Close()
		return err
	}
	switch ConnType(data[0]) {
	case CtrlConn:
		id := binary.BigEndian.Uint32(data[1:])
		dcChan := make(chan net.Conn)
		_, exist := connMap.LoadOrStore(id, dcChan)
		if exist {
			conn.Close()
			return errors.New("exist id")
		}
		go func() {
			defer connMap.Delete(id)
			timer := time.NewTimer(duration)
			select {
			case <-timer.C:
				conn.Close()
				return
			case dc := <-dcChan:
				handleFperf(conn, dc, duration)
			}
		}()
	case TransConn:
		id := binary.BigEndian.Uint32(data[1:])
		ch, ok := connMap.Load(id)
		if !ok {
			conn.Close()
			return errors.New("no ctrl conn")
		}

		ch.(chan net.Conn) <- conn
	}

	return nil
}

func handleFperf(cc, dc net.Conn, dura time.Duration) error {
	defer cc.Close()
	defer dc.Close()
	perf := NewSender(dc, cc, dura)
	// go perf.Stat.RunBandwidthIn1()
	// defer perf.Stat.StopBandwidthIn1()
	err := perf.RunSender()
	if err != nil {
		log.Println(err)
		return err
	}
	stat := perf.GetStat()
	peerStat := perf.GetPeerStat()
	tcpInfo, err := tcpinfo.GetsockoptTCPInfo(dc.(*net.TCPConn))
	if err != nil {
		log.Println(err)
		return err
	}
	r := &TcpResult{
		TCPInfo:   tcpInfo,
		SendTotal: stat.Tx,
		RecvTotal: peerStat.Rx,
		Dura:      perf.TestDuration,
		Retrans:   tcpInfo.Total_retrans * tcpInfo.Snd_mss * 8,
	}
	r.Print()
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, r.Retrans)
	_, err = cc.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

// TCPClientRecvCompatible 测试TCP性能（客户端）
func TCPClientRecvCompatible(serverAddr string, timeout int32) (r *TcpResult, err error) {
	return TCPClientRecv(serverAddr, time.Duration(timeout)*time.Second)
}

// TCPClientSend 测试TCP性能（客户端）
// func TCPClientSend(serverAddr string, duration time.Duration) (r *TcpResult, err error) {
// 	conn, err := net.Dial("tcp", serverAddr)
// 	if err != nil {
// 		log.Printf("Dial tcp failed: %s", err)
// 		return
// 	}
// 	defer conn.Close()
// 	conn.(*net.TCPConn).SetWriteBuffer(RWBufferSize)

// 	perf := NewSender(conn, duration)
// 	err = perf.RunSender()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	stat := perf.GetStat()
// 	peerStat := perf.GetPeerStat()
// 	tcpInfo, err := tcpinfo.GetsockoptTCPInfo(conn.(*net.TCPConn))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	r = &TcpResult{
// 		TCPInfo:   tcpInfo,
// 		SendTotal: stat.Tx,
// 		RecvTotal: peerStat.Rx,
// 		Dura:      duration,
// 	}
// 	return
// }

func TCPClientRecv(serverAddr string, timeout time.Duration) (r *TcpResult, err error) {
	cc, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Dial tcp failed: %s", err)
		return
	}
	defer cc.Close()
	cc.(*net.TCPConn).SetDeadline(time.Now().Add(timeout))
	id := rand.Uint32()
	data := make([]byte, FirstDataLen)
	data[0] = byte(CtrlConn)
	binary.BigEndian.PutUint32(data[1:], id)
	_, err = cc.Write(data)
	if err != nil {
		log.Printf("Send first data failed: %s", err)
		return nil, err
	}
	dc, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Dial tcp failed: %s", err)
		return
	}
	defer dc.Close()
	cc.(*net.TCPConn).SetDeadline(time.Now().Add(timeout))
	data[0] = byte(TransConn)
	_, err = dc.Write(data)
	if err != nil {
		log.Printf("Send first data failed: %s", err)
		return nil, err
	}

	perf := NewReceiver(dc, cc)
	// go perf.Stat.RunBandwidthIn1()
	// defer perf.Stat.StopBandwidthIn1()
	err = perf.RunReceiver()
	if err != nil {
		log.Println(err)
		return
	}
	stat := perf.GetStat()
	peerStat := perf.GetPeerStat()
	tcpInfo, err := tcpinfo.GetsockoptTCPInfo(dc.(*net.TCPConn))
	if err != nil {
		log.Println(err)
		return
	}
	r = &TcpResult{
		TCPInfo:   tcpInfo,
		SendTotal: peerStat.Tx,
		RecvTotal: stat.Rx,
		Dura:      perf.TestDuration,
	}
	buffer := make([]byte, 4)
	_, err = io.ReadFull(cc, buffer)
	if err != nil {
		return
	}
	r.Retrans = binary.BigEndian.Uint32(buffer)
	return
}

type TcpResult struct {
	// TCP信息
	TCPInfo *tcpinfo.TCPInfo
	// 发送的总数据量(bit)
	SendTotal uint64
	// 对端接收的总数据量(bit)
	RecvTotal uint64
	// 测试时间(second)
	Time int64
	// 重传量
	Retrans uint32
	Dura    time.Duration
}

// RetransPercents 重传率，若要以%表示需要额外*100
func (t *TcpResult) RetransPercents() float64 {
	return float64(t.Retrans) / float64(t.SendTotal)
}

// RTT 单位ms
func (t *TcpResult) RTT() int32 {
	return int32(t.TCPInfo.Rtt / 1000)
}

// Duration 测试时长（second）
func (t *TcpResult) Duration() int32 {
	return int32(t.Dura.Seconds())
}

// 发送的总数据量(bit)
func (t *TcpResult) SendTotalBit() int64 { return int64(t.SendTotal) }

// 对端接收的总数据量(bit)
func (t *TcpResult) RecvTotalBit() int64 { return int64(t.RecvTotal) }

// Print 打印测试的关键数据
func (t *TcpResult) Print() {
	fmt.Printf("Send: %d mb recv: %d mb Retran: %d(%f%%)  RTT: %d ms MSS: %d SndWin: %d\n",
		t.SendTotal/1024/1024,
		t.RecvTotal/1024/1024,
		t.Retrans,
		t.RetransPercents()*100,
		t.TCPInfo.Rtt/1000,
		t.TCPInfo.Snd_mss,
		t.TCPInfo.Snd_cwnd,
	)
	fmt.Printf("Bandwidth send: %d mbps   Peer recv: %d mbps\n",
		int64(float64(t.SendTotal)/t.Dura.Seconds()/1024/1024),
		int64(float64(t.RecvTotal)/t.Dura.Seconds()/1024/1024))
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
