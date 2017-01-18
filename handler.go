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
	if httpserver.Path(req.URL.Path).Matches(m.AuthPath) {
		reqContent := make([]byte, 1024)
		n, err := req.Body.Read(reqContent)
		if err != nil {
			return 403, err
		}

		raw, err := m.NoiseIKHandshake.Decode(reqContent[:n])
		if err != nil {
			return 403, err
		}

		back, err := m.NoiseIKHandshake.Encode([]byte(raw))
		if err != nil {
			return 403, err
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
