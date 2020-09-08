// +build !darwin

package fperf

import (
	"net"

	"github.com/brucespang/go-tcpinfo"
)

type TCPInfo struct {
	*tcpinfo.TCPInfo
}

func GetsockoptTCPInfo(tcpConn *net.TCPConn) (*TCPInfo, error) {
	i, e := tcpinfo.GetsockoptTCPInfo(tcpConn)
	return &TCPInfo{TCPInfo: i}, e
}
