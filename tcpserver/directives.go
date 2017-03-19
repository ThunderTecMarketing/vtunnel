package tcpserver

import (
	"github.com/mholt/caddy"
	"net"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/ahead"
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vtunnel/msocks"
)

func init() {

	caddy.RegisterPlugin("clients", caddy.Plugin{
		ServerType: ServerType,
		Action:     SetupClientsDirective,
	})


	caddy.RegisterPlugin("transport", caddy.Plugin{
		ServerType: ServerType,
		Action:     SetupTransportDirective,
	})
}

func SetupClientsDirective(c *caddy.Controller) (err error) {

	ctx := c.Context().(*tunnelContext)
	config := ctx.keysToConfigs[c.Key]
	var clientkey string

	if c.Next() {
		print (c.Val())
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			break
		default:
			return c.ArgErr()
		}

		for c.NextBlock() {
			userName := c.Val()
			clientkey, err = StringArg(c)
			if err != nil {
				break
			}
			config.Clients[userName] = clientkey
		}
	}

	return
}


func SetupTransportDirective(c *caddy.Controller) error {
	ctx := c.Context().(*tunnelContext)
	config := ctx.keysToConfigs[c.Key]

	config.Handler = func (ln net.Listener) error {
		context1 := &transport.TransportStreamContext2{
			Listener: ln,
		}

		context2 := new(fragment.FragmentContext)
		context3 := ahead.NewAheadContext([]byte("Key..."))
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

	return nil
}