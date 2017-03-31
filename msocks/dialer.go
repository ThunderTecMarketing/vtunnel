package msocks

import (
	"time"
	"net"
	"github.com/FTwOoO/vpncore/net/conn"
)

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

type DialerFunc func(string, string) (net.Conn, error)

func (f DialerFunc) Dial(network string, addr string) (net.Conn, error) {
	return f(network, addr)
}

type TcpDialer struct {}

func (td *TcpDialer) Dial(network, address string) (net.Conn, error) {
	return td.DialTimeout(network, address, DIAL_TIMEOUT * time.Second)
}

func (td *TcpDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

var DefaultTcpDialer Dialer = &TcpDialer{}

type ObjectDialer interface {
	Dial() (conn.ObjectIO, error)
}

