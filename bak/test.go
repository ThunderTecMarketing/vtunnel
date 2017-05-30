package main

import (
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/ahead"
	"github.com/FTwOoO/vpncore/net/conn"
	"fmt"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport"
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"github.com/FTwOoO/vtunnel/msocks"
)

func main() {

	context1 := &transport.TCPTransportStreamContext{
		Protocol:conn.PROTO_TCP,
		ListenAddr:"",
		RemoveAddr:fmt.Sprintf("127.0.0.1:%d", 10088)}

	context2 := new(fragment.FragmentContext)
	context3 := ahead.NewAheadContext([]byte("Key..."))
	context4 := new(msgpack.MsgpackContext)

	contexts := []conn.Context{context1, context2, context3, context4}
	client := new(conn.SimpleClient)

	connection, err := client.Dial(contexts)
	if err != nil {
		panic(err)
	}


	//DNS request: dig baidu.com
	packet := []byte{
		0xd3, 0x52, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x05, 0x62, 0x61, 0x69, 0x64, 0x75, 0x03, 0x63, 0x6f, 0x6d,
		0x00, 0x00, 0x01, 0x00, 0x01,
	}


	frame := &msocks.FrameDns{Data:packet}
	connection.Write(frame)

	for {
		obj, err := connection.Read()
		dnsRespFrame := obj.(*msocks.FrameDns)

		dnsResp := dnsRespFrame.Data
		dnsRespPacket := &layers.DNS{}
		err = dnsRespPacket.DecodeFromBytes(dnsResp, nil)
		if err != nil {
			print(err)
			return
		}
		print(gopacket.LayerString(dnsRespPacket))
	}
}


