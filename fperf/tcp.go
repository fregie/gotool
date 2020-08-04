package fperf

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brucespang/go-tcpinfo"
)

type TcpPerfServer struct {
	Addr string
}

func TCPServeWithContext(ctx context.Context, Addr string) error {
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
			connCh <- conn
		}
	}()
	for {
		select {
		case conn := <-connCh:
			go handleConn(conn)
		case <-ctx.Done():
			lis.Close()
			return nil
		}
	}
}

var TCPListener net.Listener

func TCPServe(Addr string) error {
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
		go handleConn(conn)
	}
}

func StopTCPServe() error {
	if TCPListener == nil {
		return errors.New("Not running")
	}
	defer func() { TCPListener = nil }()
	return TCPListener.Close()
}

func handleConn(conn net.Conn) error {
	defer conn.Close()
	perf := NewReceiver(conn)
	err := perf.RunReceiver()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// TCPClientCompatible 测试TCP性能（客户端）
func TCPClientCompatible(serverAddr string, testSeconds int32) (r *TcpResult, err error) {
	return TCPClient(serverAddr, time.Duration(testSeconds)*time.Second)
}

// TCPClient 测试TCP性能（客户端）
func TCPClient(serverAddr string, duration time.Duration) (r *TcpResult, err error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Dial tcp failed: %s", err)
		return
	}
	defer conn.Close()

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

type TcpResult struct {
	// TCP信息
	TCPInfo *tcpinfo.TCPInfo
	// 发送的总数据量(bit)
	SendTotal uint64
	// 对端接收的总数据量(bit)
	RecvTotal uint64
	// 测试时间(second)
	Time int64
	Dura time.Duration
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

// Print 打印测试的关键数据
func (t *TcpResult) Print() {
	fmt.Printf("Send: %d mb  Peer recv: %d mb Retran: %d(%f%%)  RTT: %d ms\n",
		t.SendTotal/1024/1024,
		t.RecvTotal/1024/1024,
		t.TCPInfo.Total_retrans,
		float64(t.TCPInfo.Total_retrans)*float64(t.TCPInfo.Snd_mss)*100/float64(t.SendTotal),
		t.TCPInfo.Rtt/1000,
	)
	fmt.Printf("Bandwidth send: %d mbps   Peer recv: %d mbps\n", t.SendTotal/uint64(t.Dura.Seconds())/1024/1024, t.RecvTotal/uint64(t.Dura.Seconds())/1024/1024)
}
