package tcpclient

import (
	"github.com/ginuerzh/gosocks5"
	"net"
	"io"
	"github.com/FTwOoO/vtunnel/msocks"
)

type NoAuthSocksServerSelector struct{}

func (selector *NoAuthSocksServerSelector) Methods() []uint8 {
	return []uint8{gosocks5.MethodNoAuth}
}

func (selector *NoAuthSocksServerSelector) Select(methods ...uint8) (method uint8) {

	method = gosocks5.MethodNoAcceptable
	for _, m := range methods {
		if m == gosocks5.MethodNoAuth {
			return gosocks5.MethodNoAuth
		}
	}
	return
}

func (selector *NoAuthSocksServerSelector) OnSelected(method uint8, conn net.Conn) (net.Conn, error) {

	switch method {
	case gosocks5.MethodNoAcceptable:
		return nil, gosocks5.ErrBadMethod
	}

	return conn, nil
}

type Socks5Server struct {
	Socks5ListenAddr string
	Selector         gosocks5.Selector
	Pool             *msocks.SessionPool
}

func (s *Socks5Server) Serve() error {
	ln, err := net.Listen("tcp", s.Socks5ListenAddr)

	if err != nil {
		return nil
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Socks5Server) handleConn(conn net.Conn) {
	defer conn.Close()

	conn = gosocks5.ServerConn(conn, s.Selector)
	req, err := gosocks5.ReadRequest(conn)
	if err != nil {
		return
	}

	s.HandleRequest(conn, req)
}

func (s *Socks5Server) HandleRequest(conn net.Conn, req *gosocks5.Request) (err error) {

	switch req.Cmd {
	case gosocks5.CmdConnect:
		s.handleConnect(conn, req)
	default:
		rep := gosocks5.NewReply(gosocks5.CmdUnsupported, nil)
		if err = rep.Write(conn); err != nil {
			return
		}
	}

	return
}

func (s *Socks5Server) handleConnect(conn net.Conn, req *gosocks5.Request) {
	session, err := s.Pool.Get()
	if session != nil {
		rep := gosocks5.NewReply(gosocks5.NetUnreachable, nil)
		rep.Write(conn)
		return
	}

	cc, err := session.Dial("tcp", req.Addr.String())
	if err != nil {
		rep := gosocks5.NewReply(gosocks5.HostUnreachable, nil)
		rep.Write(conn)
		return
	} else {
		defer cc.Close()

		rep := gosocks5.NewReply(gosocks5.Succeeded, nil)
		if err = rep.Write(conn); err != nil {
			return
		}
		s.connected(cc, conn)

	}
	return
}

func (s *Socks5Server) connected(conn1, conn2 net.Conn) (err error) {
	errc := make(chan error, 2)

	go func() {
		_, err := io.Copy(conn1, conn2)
		errc <- err
	}()

	go func() {
		_, err := io.Copy(conn2, conn1)
		errc <- err
	}()

	select {
	case err = <-errc:
	}

	return
}

