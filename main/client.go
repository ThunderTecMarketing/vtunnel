package main

import (
	"github.com/FTwOoO/vtunnel/tcpclient"
	"github.com/FTwOoO/vtunnel/msocks"
)

func main() {
	server := &tcpclient.Socks5Server{
		Socks5ListenAddr: "0.0.0.0:10808",
		Selector:new(tcpclient.NoAuthSocksServerSelector),
		Pool:msocks.CreateSessionPool(0,0,[]msocks.ObjectDialer{&tcpclient.ClientDialer{RemoteAddr:"127.0.0.1:10088", Key:"Key..."}}),
	}
	server.Serve()
}
