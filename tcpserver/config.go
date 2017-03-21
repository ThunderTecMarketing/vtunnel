/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package tcpserver

import (
	"time"
	"errors"
	"net"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/ahead"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport"
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"github.com/FTwOoO/vtunnel/msocks"
	"github.com/FTwOoO/vtunnel/tcpclient"
)

var DefaultPrologue = "vtunnel"
var KeyLength = 16

var DefaultPeerTimeout = time.Duration(30 * time.Second)
var DefaultTokenTimeout = time.Duration(10 * time.Second)

var ErrInValidHandshakeStep = errors.New("Invalid handshake step")
var ErrInValidKeyLength = errors.New("Invalid key length")
var ErrPeerAlreadyExist = errors.New("Peer exists")
var ErrPacketLengthInvalid = errors.New("Packet length is not in (0,MTU]")
var ErrWriteFail = errors.New("Write fail")

type TransportType string

var (
	TRANSPORT1 = TransportType("TCP-Fragment-AheadGCM256-Msgpack")
)

type ListenerHandler func(net.Listener) error

type ServerConfig struct {
	ListenAddr     string
	RemoteAddr     string

	IsServer     bool
	TransportKey   []byte
	TransportType  TransportType
	LocalProxyType string
}

func (s *ServerConfig) GetHandler() ListenerHandler {
	if s.IsServer == true && s.TransportType == TRANSPORT1 {
		return func(ln net.Listener) error {
			context1 := &transport.TransportStreamContext2{
				Listener: ln,
			}

			context2 := new(fragment.FragmentContext)
			context3 := ahead.NewAheadContext([]byte(s.TransportKey))
			context4 := new(msgpack.MsgpackContext)

			contexts := []conn.Context{context1, context2, context3, context4}
			server := new(conn.SimpleServer)
			listener, err := server.NewListener(contexts)
			if err != nil {
				return err
			}

			msocksServer, err := msocks.NewMsocksServer(msocks.DefaultTcpDialer)
			if err != nil {
				return err
			}

			return msocksServer.Serve(listener)

		}
	}
	if s.IsServer == false && s.TransportType == TRANSPORT1 {
		return func(ln net.Listener) error {
			dialer := &Dialer{
				Pool:msocks.CreateSessionPool(0, 0,
					[]msocks.ObjectDialer{&tcpclient.ClientDialer{
						RemoteAddr:s.RemoteAddr,
						Key:s.TransportKey},
					}),
			}

			server := &tcpclient.Socks5Server{
				Socks5ListenAddr: "",
				Selector:new(tcpclient.NoAuthSocksServerSelector),
				Dialer:dialer,
			}
			return server.Serve(ln)
		}
	}

	return nil
}

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