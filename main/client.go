package main

import (
	"github.com/FTwOoO/vtunnel/tcpclient"
	"github.com/FTwOoO/vtunnel/msocks"
	"net"
)

type Dialer struct {
	Pool *msocks.SessionPool
}

func (d *Dialer) Dial(network string, addr string) (net.Conn, error) {
	session, err := d.Pool.Get()
	if err != nil {
		return nil, err
	}
	return session.Dial(network, addr)
}

func main() {
	dialer := &Dialer{Pool:msocks.CreateSessionPool(0, 0,
		[]msocks.ObjectDialer{&tcpclient.ClientDialer{
			RemoteAddr:"127.0.0.1:10088",
			Key:"Key..."},
		}),
	}

	server := &tcpclient.Socks5Server{
		Socks5ListenAddr: "0.0.0.0:10808",
		Selector:new(tcpclient.NoAuthSocksServerSelector),
		Dialer:dialer,
	}
	server.Serve()
}
