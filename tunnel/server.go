/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package tunnel

import (
	"sync"
	"time"
	"net"
	"runtime"
	"github.com/FTwOoO/vtunnel/socks5_server"
	"github.com/FTwOoO/vtunnel/config"
	"github.com/FTwOoO/vtunnel/proxy_dialer"
)

type Server struct {
	listener   net.Listener
	listenerMu sync.Mutex
	config     *config.Config

	connTimeout time.Duration // max time to wait for a connection before force stop

	doneChan chan struct{}
}

var GracefulTimeout = 5 * time.Second

func NewServer(config *config.Config) (*Server, error) {
	s := &Server{
		config:      config,
		connTimeout: GracefulTimeout,
		doneChan:    make(chan struct{}),
	}

	return s, nil
}

func (s *Server) Listen() (net.Listener, error) {

	ln, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		var succeeded bool
		if runtime.GOOS == "windows" {
			// Windows has been known to keep sockets open even after closing the listeners.
			// Tests reveal this error case easily because they call Start() then Stop()
			// in succession. TODO: Better way to handle this? And why limit this to Windows?
			for i := 0; i < 20; i++ {
				time.Sleep(100 * time.Millisecond)
				ln, err = net.Listen("tcp", s.config.ListenAddr)
				if err == nil {
					succeeded = true
					break
				}
			}
		}
		if !succeeded {
			return nil, err
		}
	}

	if tcpLn, ok := ln.(*net.TCPListener); ok {
		ln = tcpKeepAliveListener{TCPListener: tcpLn}
	}

	// Very important to return a concrete caddy.Listener
	// implementation for graceful restarts.
	return ln, nil
}

func (s *Server) Serve(ln net.Listener) error {
	s.listenerMu.Lock()
	s.listener = ln
	s.listenerMu.Unlock()

	diler := &proxy_dialer.HttpConnectDialer{C: s.config}

	sock5Serv := socks5_server.Socks5Server{
		Selector: new(socks5_server.NoAuthSocksServerSelector),
		Dialer:   diler,
	}

	return sock5Serv.Serve(ln)

}

func (s *Server) Address() string {
	return s.config.ListenAddr
}

func (s *Server) Stop() error {
	close(s.doneChan)
	return nil
}
