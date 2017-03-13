package tcpserver

import "io"

type Handler interface {
	ServeHTTP(io.Writer) (int, error)
}

type handler struct {
	Next      Handler
	DnsServer *DnsServer
}
