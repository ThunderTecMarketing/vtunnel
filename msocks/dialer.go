package msocks

import (
	"time"
	"net"
	"github.com/FTwOoO/vpncore/net/conn"
)

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}


type TcpDialer struct {
}

func (td *TcpDialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

func (td *TcpDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

var DefaultTcpDialer Dialer = &TcpDialer{}

type ObjectDialer interface {
	Dial() (conn.ObjectIO, error)
}

