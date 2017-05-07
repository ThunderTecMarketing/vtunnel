package client

import (
	"net"
	"github.com/FTwOoO/vtunnel/msocks"
	"github.com/FTwOoO/vtunnel/tunnel"
)

func init() {
	tunnel.ResgisterHandlerGenerator(GetHandler)
}

func GetHandler(config *tunnel.Config) tunnel.ListenerHandler {
	if config.IsServer == false && config.TransportType == tunnel.TRANSPORT1 {
		return func(ln net.Listener) error {

			dialer := &MultiplexerDialer{
				Pool:msocks.NewClientSessionPool(0, 30,
					[]msocks.ObjectDialer{&ProtocolDialer{
						RemoteAddr:config.RemoteAddr,
						Key:config.TransportKey},
					}),
			}

			s := &Socks5Server{
				Selector:new(NoAuthSocksServerSelector),
				Dialer:dialer,
			}

			return s.Serve(ln)
		}
	}

	return nil
}
