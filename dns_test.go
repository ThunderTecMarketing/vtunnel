/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-02
 */

package vpn

import (
	"testing"
	"fmt"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"net/http"
	"errors"
	"net/http/httptest"
	"bytes"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
)

func createTestHandler() *handler {

	input := fmt.Sprintf(`vpn /vpn {
		    publickey %s
		    privatekey %s
		    clients {
		 	%s
		    }

		    subnet 192.168.4.1/24
		    mtu 1400
		    dnsport 53
		}`, serverPublicKey, serverPrivateKey, validClientPublicKey)

	h, err := Parse(caddy.NewTestController("http", input))
	if err != nil {
		print(err)
		return nil
	}

	h.Next = httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusNotFound, errors.New("404 error")
	})

	return h
}

func TestDNS(t *testing.T) {
	h := createTestHandler()

	//DNS request: dig baidu.com
	packet := []byte{
		0xd3, 0x52, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x05, 0x62, 0x61, 0x69, 0x64, 0x75, 0x03, 0x63, 0x6f, 0x6d,
		0x00, 0x00, 0x01, 0x00, 0x01,
	}
	buf := bytes.NewBuffer(packet)


	req, err := http.NewRequest("POST", "https://127.0.0.1/dns/", buf)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	rec := httptest.NewRecorder()
	statusCode, err := h.ServeHTTP(rec, req)

	if statusCode != http.StatusOK {
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	dnsResp := make([]byte, 1024)
	n, err := rec.Body.Read(dnsResp)

	dnsRespPacket := &layers.DNS{}
	err = dnsRespPacket.DecodeFromBytes(dnsResp[:n], nil)
	if err != nil {
		t.Fatalf("Could not parse DNS response: %v", err)
	}
	print(gopacket.LayerString(dnsRespPacket))

}