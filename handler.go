package vpn

import (
	"net"
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type handler struct {
	Next             httpserver.Handler

	PublicKey        string
	PrivateKey       string
	ClientPublicKeys []string
	Ip               net.IP
	Subnet           *net.IPNet
	MTU              uint16
	DnsPort          uint16
	AuthPath         string
	PacketPath       string
}

func (m *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) (int, error) {
	// if the request path is any of the configured paths
	// write hello

	if httpserver.Path(req.URL.Path).Matches(m.AuthPath) {
		w.Write([]byte("Auth OK!"))
		return 200, nil
	}

	if httpserver.Path(req.URL.Path).Matches(m.PacketPath) {
		w.Write([]byte("Packet OK!"))
		return 200, nil
	}

	return m.Next.ServeHTTP(w, req)
}
