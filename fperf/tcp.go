package fperf

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/brucespang/go-tcpinfo"
)

const million = 1048576

type TcpPerfServer struct {
	Addr string
}

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
			conn.(*net.TCPConn).SetWriteBuffer(million)
			connCh <- conn
		}
	}()
	for {
		select {
		case conn := <-connCh:
			go handleConn(conn, duration)
		case <-ctx.Done():
			lis.Close()
			return nil
		}
	}
}

var TCPListener net.Listener

func TCPSendServe(Addr string, duration time.Duration) error {
	if TCPListener != nil {
		return errors.New("already running")
	}
	var err error
	TCPListener, err = net.Listen("tcp", Addr)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s", TCPListener.Addr().String())
	for {
		conn, err := TCPListener.Accept()
		if err != nil {
			log.Printf("accept err: %s", err)
			return err
		}
		conn.(*net.TCPConn).SetWriteBuffer(million)
		go handleConn(conn, duration)
	}
}

func StopTCPServe() error {
	if TCPListener == nil {
		return errors.New("Not running")
	}
	defer func() { TCPListener = nil }()
	return TCPListener.Close()
}

func handleConn(conn net.Conn, dura time.Duration) error {
	defer conn.Close()
	perf := NewSender(conn, dura)
	// go perf.Stat.RunBandwidthIn1()
	// defer perf.Stat.StopBandwidthIn1()
	err := perf.RunSender()
	if err != nil {
		log.Println(err)
		return err
	}
	stat := perf.GetStat()
	peerStat := perf.GetPeerStat()
	tcpInfo, err := tcpinfo.GetsockoptTCPInfo(conn.(*net.TCPConn))
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
	_, err = conn.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

// TCPClientSendCompatible 测试TCP性能（客户端）
func TCPClientSendCompatible(serverAddr string, testSeconds int32) (r *TcpResult, err error) {
	return TCPClientSend(serverAddr, time.Duration(testSeconds)*time.Second)
}

// TCPClientRecvCompatible 测试TCP性能（客户端）
func TCPClientRecvCompatible(serverAddr string, testSeconds int32) (r *TcpResult, err error) {
	return TCPClientRecv(serverAddr)
}

// TCPClientSend 测试TCP性能（客户端）
func TCPClientSend(serverAddr string, duration time.Duration) (r *TcpResult, err error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Dial tcp failed: %s", err)
		return
	}
	defer conn.Close()
	conn.(*net.TCPConn).SetWriteBuffer(million)

	perf := NewSender(conn, duration)
	err = perf.RunSender()
	if err != nil {
		log.Println(err)
		return
	}
	stat := perf.GetStat()
	peerStat := perf.GetPeerStat()
	tcpInfo, err := tcpinfo.GetsockoptTCPInfo(conn.(*net.TCPConn))
	if err != nil {
		log.Println(err)
		return
	}
	r = &TcpResult{
		TCPInfo:   tcpInfo,
		SendTotal: stat.Tx,
		RecvTotal: peerStat.Rx,
		Dura:      duration,
	}
	return
}

func TCPClientRecv(serverAddr string) (r *TcpResult, err error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Dial tcp failed: %s", err)
		return
	}
	defer conn.Close()
	conn.(*net.TCPConn).SetReadBuffer(million)

	perf := NewReceiver(conn)
	// go perf.Stat.RunBandwidthIn1()
	// defer perf.Stat.StopBandwidthIn1()
	err = perf.RunReceiver()
	if err != nil {
		log.Println(err)
		return
	}
	stat := perf.GetStat()
	peerStat := perf.GetPeerStat()
	tcpInfo, err := tcpinfo.GetsockoptTCPInfo(conn.(*net.TCPConn))
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
	_, err = io.ReadFull(conn, buffer)
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
	return float64(t.TCPInfo.Total_retrans) * float64(t.TCPInfo.Snd_mss) / float64(t.SendTotal)
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
	fmt.Printf("Send: %d mb recv: %d mb Retran: %d(%f%%)  RTT: %d ms\n",
		t.SendTotal/1024/1024,
		t.RecvTotal/1024/1024,
		t.Retrans,
		float64(t.Retrans)*100/float64(t.SendTotal),
		t.TCPInfo.Rtt/1000,
	)
	fmt.Printf("Bandwidth send: %d mbps   Peer recv: %d mbps\n",
		int64(float64(t.SendTotal)/t.Dura.Seconds()/1024/1024),
		int64(float64(t.RecvTotal)/t.Dura.Seconds()/1024/1024))
}
