package vpn

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/athom/goset"
	"errors"

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
			return http.StatusUnauthorized, err
		}

		rawContent, err := m.NoiseIKHandshake.Decode(reqContent[:n])
		if err != nil {
			return http.StatusUnauthorized, err
		}

		rs := m.NoiseIKHandshake.Hs.PeerStatic()
		if !goset.IsIncluded(m.Config.ClientPublicKeys, rs) {
			return http.StatusUnauthorized, errors.New("Invalid Key")
		}

		respContent, err := m.NoiseIKHandshake.Encode([]byte(rawContent))
		if err != nil {
			return http.StatusUnauthorized, err
		}

		w.Write(respContent)
		return http.StatusOK, nil
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

		return http.StatusOK, nil
	}

	return m.Next.ServeHTTP(w, req)
}
