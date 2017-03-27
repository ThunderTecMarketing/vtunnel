package server

import (
	"net"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/ahead"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport"
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"github.com/FTwOoO/vtunnel/msocks"
	"github.com/FTwOoO/vtunnel/tunnel"
)

func init() {
	tunnel.ResgisterHandlerGenerator(GetHandler)
}

func GetHandler(config *tunnel.Config) tunnel.ListenerHandler {
	if config.IsServer == true && config.TransportType == tunnel.TRANSPORT1 {
		return func(ln net.Listener) error {
			context1 := &transport.TransportStreamContext2{
				Listener: ln,
			}

			context2 := new(fragment.FragmentContext)
			context3 := ahead.NewAheadContext([]byte(config.TransportKey))
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

	return nil
}