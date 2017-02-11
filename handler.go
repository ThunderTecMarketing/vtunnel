package vpn

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"errors"
	"net"
	"github.com/athom/goset"
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

	if httpserver.Path(req.URL.Path).Matches(m.VPNAuthPath) {

		var err error
		var password []byte
		var peer *Peer

		_, passwordS, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"vpn\"")
			return http.StatusUnauthorized, nil
		}

		password, err = hex.DecodeString(passwordS)
		if err != nil {
			return http.StatusUnauthorized, err
		}

		if !goset.IsIncluded(m.Config.ClientPublicKeys, password) {
			return http.StatusUnauthorized, errors.New("Invalid Key")
		}

		peer, err = m.Peers.AddPeer(password)
		if err != nil {
			return http.StatusUnauthorized, err
		}

		clientSetting := ClientSetting{Ip:peer.Ip, Subnet:m.Subnet, Mtu:m.MTU}
		var respContent []byte = []byte(clientSetting.Encode())
		if err != nil {
			return http.StatusUnauthorized, err
		}

		w.Write(respContent)
		return http.StatusOK, nil
	}

	if httpserver.Path(req.URL.Path).Matches(m.VPNDataPath) {
		var err error
		var password []byte
		var peer *Peer

		ipS, passwordS, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"vpn\"")
			return http.StatusUnauthorized, nil
		}

		password, err = hex.DecodeString(passwordS)
		if err != nil {
			return http.StatusUnauthorized, err
		}

		ip := net.ParseIP(ipS)
		if ip == nil {
			return http.StatusUnauthorized, errors.New("Ip format error")
		}

		// remove credentials from request to avoid leaking upstream
		req.Header.Del("Authorization")

		peer = m.Peers.GetPeerByIp(ip)
		if peer == nil || !bytes.Equal(peer.Key, password) {
			return http.StatusUnauthorized, errors.New("Invalid ip/password ")
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