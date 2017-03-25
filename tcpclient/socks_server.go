package tcpclient

import (
	"github.com/ginuerzh/gosocks5"
	"net"
	"io"
	"github.com/FTwOoO/vtunnel/msocks"
	"github.com/pkg/errors"
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
	Selector         gosocks5.Selector
	Dialer           msocks.Dialer
}

func (s *Socks5Server) Serve(ln interface{}) (err error) {

	var listener net.Listener
	switch ln.(type) {
	case net.Listener:
		listener = ln.(net.Listener)

	case string:
		socks5ListenAddr := ln.(string)
		listener, err = net.Listen("tcp", socks5ListenAddr)
		if err != nil {
			return err
		}
	default:
		return errors.New("Unkown type")
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}

	return nil
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
	cc, err := s.Dialer.Dial("tcp", req.Addr.String())
	if err != nil {
		rep := gosocks5.NewReply(gosocks5.NetUnreachable, nil)
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

