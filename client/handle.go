package client

import (
	"net"
	"github.com/FTwOoO/vtunnel/msocks"
	"github.com/FTwOoO/vtunnel/tunnel"
	"github.com/FTwOoO/vtunnel/util"
	"github.com/FTwOoO/vpncore/net/gfw"
)

func init() {
	tunnel.ResgisterHandlerGenerator(GetHandler)
}

func GetHandler(config *tunnel.Config) tunnel.ListenerHandler {
	if config.IsServer == false && config.TransportType == tunnel.TRANSPORT1 {
		return func(ln net.Listener) error {

			var gfwlist *gfw.ItemSet = nil
			filePath, err := util.GetConfigPath(config.GFWListFile)
			if err == nil {
				gfwlist = gfw.NewItemSet(filePath, 5000)
				gfwlist.Load()
			} else {
				gfwlist, err = gfw.CreateGFWList("", config.GFWListFile)
				if err != nil {
					return err
				}
			}

			dialer := &GFWDialer{
				Gfwlist:gfwlist,
				Pool:msocks.CreateSessionPool(0, 30,
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
