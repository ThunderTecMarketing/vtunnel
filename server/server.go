package server

import (
	"errors"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vpncore/net/dns"
	"github.com/FTwOoO/vtunnel/msocks"
)

type Server struct {
	remoteDialer msocks.Dialer
	dnsServer    *dns.DNSServer
}

func NewServer(dialer msocks.Dialer) (ms *Server, err error) {
	if dialer == nil {
		err = errors.New("empty dialer")
		return
	}

	dnsServer, err := dns.NewDNSServer(nil, true)
	if err != nil {
		return
	}

	ms = &Server{
		remoteDialer: dialer,
		dnsServer: dnsServer,
	}

	return
}

func (ms *Server) Handler(conn conn.ObjectIO) {
	//log.Notice("connection come from: %s => %s.", conn.RemoteAddr(), conn.LocalAddr())

	sess := msocks.NewSession(conn, ms.dnsServer, ms.remoteDialer)
	sess.Run()

	//log.Notice("server session %d quit: %s => %s.", sess.LocalPort(), conn.RemoteAddr(), conn.LocalAddr())
}

func (ms *Server) Serve(listener conn.ObjectListener) (err error) {
	var connection conn.ObjectIO

	for {
		connection, err = listener.Accept()
		if err != nil {
			continue
		}
		go func() {
			defer connection.Close()
			ms.Handler(connection)
		}()
	}
	return
}
