// Author: xiexu
// Date: 2022-09-20

//go:build !windows

// github.com/ixugo/netpulse
// 有更多关于 ip 的处理
package system

import (
	"net"
	"strconv"
	"strings"
)

// PortUsed 检测端口  true:已使用;false:未使用
func PortUsed(mode string, port int) bool {
	if port > 65535 || port < 0 {
		return true
	}

	switch strings.ToLower(mode) {
	case "tcp":
		return TCPPortUsed(port)
	default:
		return UDPPortUsed(port)
	}
}

// TCPPortUsed 通过尝试监听指定 TCP 端口来检测其是否已被占用
func TCPPortUsed(port int) bool {
	addr, _ := net.ResolveTCPAddr("tcp", net.JoinHostPort("", strconv.Itoa(port)))
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return true
	}
	_ = conn.Close()
	return false
}

// UDPPortUsed 通过尝试监听指定 UDP 端口来检测其是否已被占用
func UDPPortUsed(port int) bool {
	addr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort("", strconv.Itoa(port)))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return true
	}
	_ = conn.Close()
	return false
}

// Deprecated: 使用 github.com/ixugo/netpulse/ip 替代
func ExternalIP() (string, error) {
	panic("deprecated")
}

// Deprecated: 使用 github.com/ixugo/netpulse/ip 替代
// ip.InternalIP()
func LocalIP() string {
	panic("deprecated")
}
