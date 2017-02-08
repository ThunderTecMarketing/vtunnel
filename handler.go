package vpn

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/FTwOoO/noise"
	"github.com/athom/goset"
	"errors"
	"net"
	"encoding/hex"
	"bytes"
)

const useNoiseIX = false

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

		var ixHandshake *NoiseIXHandshake
		var rs []byte

		if useNoiseIX {
			ixHandshake, err = NewNoiseIXHandshake(
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

			rs = ixHandshake.Hs.PeerStatic()

		} else {
			ixHandshake = nil
			rs = reqContent[:n]
		}

		if !goset.IsIncluded(m.Config.ClientPublicKeys, rs) {
			return http.StatusUnauthorized, errors.New("Invalid Key")
		}
		peer, err := m.Peers.AddPeer(rs, ixHandshake, NewToken(DefaultTokenTimeout))
		if err != nil {
			return http.StatusUnauthorized, err
		}

		clientSetting := ClientSetting{Ip:peer.Ip, Subnet:m.Subnet, Mtu:m.MTU, Token:peer.Token.Value}
		var respContent []byte
		if useNoiseIX {
			respContent, err = ixHandshake.Encode([]byte(clientSetting.Encode()))
			if err != nil {
				return http.StatusUnauthorized, err
			}
		} else {
			respContent = []byte(clientSetting.Encode())
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

		var peer *Peer
		if peer = m.Peers.GetPeerByIp(ip); peer == nil || !peer.IsValid() || !bytes.Equal(peer.Token.Value, token) {
			return http.StatusUnauthorized, errors.New("Invalid token or peer ")
		}

		packets, err := ReadPackets(req.Body)
		if err != nil {
			return http.StatusBadRequest, err
		}

		for _, packet := range packets {
			m.Fowarder.Send(packet)
		}

		//TODO: limit the num of packets per resp
		packetsToWrite := m.Fowarder.Recv(peer.Ip)
		if len(packetsToWrite) > 0 {
			err = WritePackets(w, packetsToWrite)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

		return http.StatusOK, nil
	}

	return m.Next.ServeHTTP(w, req)
}

func (m *handler) DeletePeer(peer *Peer) {
	m.Fowarder.DeleteTarget(peer.Ip)
}