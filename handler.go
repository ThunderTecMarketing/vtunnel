package vpn

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type handler struct {
	Config
	Next             httpserver.Handler
	Fowarder         *Fowarder
	NoiseIKHandshake *NoiseIXHandshake
}

func (m *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) (int, error) {
	// if the request path is any of the configured paths
	// write hello

	if httpserver.Path(req.URL.Path).Matches(m.AuthPath) {
		reqContent := make([]byte, 1024)
		n, err := req.Body.Read(reqContent)
		if err != nil {
			return 403, nil
		}

		raw, err := m.NoiseIKHandshake.Decode(reqContent[:n])
		if err != nil {
			return 403, nil
		}

		back, err := m.NoiseIKHandshake.Encode([]byte(raw))
		if err != nil {
			return 403, nil
		}

		w.Write(back)
		return 200, nil
	}

	if httpserver.Path(req.URL.Path).Matches(m.PacketPath) {
		newPacket := make([]byte, m.Config.MTU)
		n, err := req.Body.Read(newPacket)
		if err == nil {
			m.Fowarder.Send(newPacket[:n])
		}

		_, err = m.Fowarder.Recv()
		if err != nil {
			w.Write([]byte{})
		}

		return 200, nil
	}

	return m.Next.ServeHTTP(w, req)
}
