package vpn

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type handler struct {
	Config
	Next httpserver.Handler
	Fowarder Fowarder
}

func (m *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) (int, error) {
	// if the request path is any of the configured paths
	// write hello

	if httpserver.Path(req.URL.Path).Matches(m.AuthPath) {
		w.Write([]byte("Auth OK!"))
		return 200, nil
	}

	if httpserver.Path(req.URL.Path).Matches(m.PacketPath) {
		newPacket := make([]byte, m.Config.MTU)
		n, err := req.Body.Read(newPacket)
		if err == nil {
			m.Fowarder.Send(newPacket[:n])
		}

		views, err := m.Fowarder.Recv()
		if err != nil {
			w.Write([]byte{})
		}

		return 200, nil
	}

	return m.Next.ServeHTTP(w, req)
}
