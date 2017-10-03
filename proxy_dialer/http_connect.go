package proxy_dialer

import (
	"bufio"
	"net/http"
	"crypto/tls"
	"net"
	"golang.org/x/net/proxy"
	"errors"
	"strings"
	"github.com/FTwOoO/vtunnel/config"
	"fmt"
)

func dial(proxyAddr string, useTls bool) (net.Conn, error) {
	if useTls {
		return tls.Dial("tcp", proxyAddr, &tls.Config{InsecureSkipVerify: true})
	} else {
		return net.Dial("tcp", proxyAddr)
	}
}

var _ proxy.Dialer = new(HttpConnectDialer)

type HttpConnectDialer struct {
	C *config.Config
}

func (d *HttpConnectDialer) Dial(network, addr string) (proxyConn net.Conn, err error) {

	if network != "tcp" {
		return nil, errors.New("Unsupported network")
	}

	s := strings.Split(addr, ":")
	hostAddr, _ := s[0], s[1]

	proxyConn, err = dial(d.C.ProxyServer.ProxyAddr, true)
	if err != nil {
		return
	}

	defer func() {
		if err != nil && proxyConn != nil {
			proxyConn.Close()
		}
	}()

	connectRequest := http.Request{Header: make(http.Header)}
	if len(d.C.ProxyServer.BasicProxyCredentials) > 0 {
		connectRequest.Header.Set("Proxy-Authorization", d.C.ProxyServer.BasicProxyCredentials)
	}

	connectRequest.Host = hostAddr
	connectRequest.RequestURI = addr
	connectRequest.Method = "CONNECT"

	switch d.C.ProxyServer.HTTPVer {
	case "HTTP/2.0":
		connectRequest.ProtoMajor = 2
		connectRequest.ProtoMinor = 0
	case "HTTP/1.1":
		connectRequest.ProtoMajor = 1
		connectRequest.ProtoMinor = 1
	default:
		panic("http2ProxyVer: " + d.C.ProxyServer.HTTPVer)
	}
	connectRequest.Proto = d.C.ProxyServer.HTTPVer

	err = connectRequest.Write(proxyConn)
	if err != nil {
		return
	}
	connectResponse, err2 := http.ReadResponse(bufio.NewReader(proxyConn), &connectRequest)
	if err2 != nil {
		return nil, err2
	} else if connectResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CONNECT fail: HTTP[%d], %v", connectResponse.StatusCode, connectResponse)
	}

	return
}
