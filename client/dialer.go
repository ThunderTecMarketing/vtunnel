package client

import (
	"net"
	"github.com/FTwOoO/vtunnel/msocks"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/ahead"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport"
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"github.com/FTwOoO/vpncore/net/gfw"
)

type ContextDialer interface {
	Dial(srcAddr net.Addr, network string, addr string) (net.Conn, error)
}

//implements ContextDialer
type GFWDialer struct {
	Pool    *msocks.SessionPool
	Gfwlist *gfw.ItemSet
}

func (d *GFWDialer) Dial(srcAddr net.Addr, network string, addr string) (net.Conn, error) {

	if d.Gfwlist.Hit(addr) {
		return net.Dial(network, addr)
	}

	session, err := d.Pool.Get()
	if err != nil {
		return nil, err
	}

	return session.Dial(srcAddr, network, addr)
}

type ProtocolDialer struct {
	RemoteAddr string
	Key        []byte
}

func (c *ProtocolDialer) Dial() (connection conn.ObjectIO, err error) {
	context1 := &transport.TransportStreamContext{
		Protocol:conn.PROTO_TCP,
		ListenAddr:"",
		RemoveAddr:c.RemoteAddr}

	context2 := new(fragment.FragmentContext)
	context3 := ahead.NewAheadContext([]byte(c.Key))
	context4 := new(msgpack.MsgpackContext)

	contexts := []conn.Context{context1, context2, context3, context4}
	client := new(conn.SimpleClient)
	connection, err = client.Dial(contexts)
	return
}
