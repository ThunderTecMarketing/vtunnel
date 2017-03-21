package tcpclient

import (
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/ahead"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport"
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
)

type ClientDialer struct {
	RemoteAddr string
	Key        []byte
}

func (c *ClientDialer) Dial() (connection conn.ObjectIO, err error) {
	context1 := &transport.TransportStreamContext{
		Protocol:conn.PROTO_TCP,
		ListenAddr:"",
		RemoveAddr:c.RemoteAddr}

	context2 := new(fragment.FragmentContext)
	context3 := ahead.NewAheadContext([]byte(c.Key))  //"Key..."))
	context4 := new(msgpack.MsgpackContext)

	contexts := []conn.Context{context1, context2, context3, context4}
	client := new(conn.SimpleClient)
	connection, err = client.Dial(contexts)
	return
}
