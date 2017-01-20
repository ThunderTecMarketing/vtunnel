package vpn

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/athom/goset"
	"errors"
	"github.com/FTwOoO/noise"
	"net"
	"encoding/hex"
	"bytes"
)

type handler struct {
	Config
	Next     httpserver.Handler
	Fowarder *Fowarder
	Peers    *Peers
}

func (m *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) (int, error) {
	if httpserver.Path(req.URL.Path).Matches(m.AuthPath) {
		reqContent := make([]byte, 1024)
		n, err := req.Body.Read(reqContent)
		if err != nil {
			return http.StatusUnauthorized, err
		}

		ixHandshake, err := NewNoiseIXHandshake(
			DefaultCipherSuite,
			[]byte(DefaultPrologue),
			noise.DHKey{Public:m.PublicKey, Private:m.PrivateKey},
			false,
		)

		if n <= 0 {
			return http.StatusBadRequest, errors.New("Need HTTP body")
		}

		_, err = ixHandshake.Decode(reqContent[:n])
		if err != nil {
			return http.StatusUnauthorized, err
		}

		rs := ixHandshake.Hs.PeerStatic()
		if !goset.IsIncluded(m.Config.ClientPublicKeys, rs) {
			return http.StatusUnauthorized, errors.New("Invalid Key")
		}

		peer, err := m.Peers.AddPeer(rs, ixHandshake, NewToken(DefaultTokenTimeout))
		if err != nil {
			return http.StatusUnauthorized, err
		}

		clientSetting:=ClientSetting{Ip:peer.Ip, Subnet:m.Subnet, Mtu:m.MTU, Token:peer.Token.Value}
		respContent, err := ixHandshake.Encode([]byte(clientSetting.Encode()))
		if err != nil {
			return http.StatusUnauthorized, err
		}

		w.Write(respContent)
		return http.StatusOK, nil
	}

	if httpserver.Path(req.URL.Path).Matches(m.PacketPath) {
		ipS, tokenV, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Restricted\"")
			return http.StatusUnauthorized, nil
		}

		// remove credentials from request to avoid leaking upstream
		req.Header.Del("Authorization")

		ip := net.ParseIP(ipS)
		if ip == nil {
			return http.StatusUnauthorized, errors.New("Ip format error")
		}

		token, err := hex.DecodeString(tokenV)
		if err != nil {
			return http.StatusUnauthorized, err
		}

		if peer := m.Peers.GetPeerByIp(ip); peer == nil || !peer.IsValid() || !bytes.Equal(peer.Token.Value, token) {
			return http.StatusUnauthorized, errors.New("Invalid token or peer ")
		}

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

