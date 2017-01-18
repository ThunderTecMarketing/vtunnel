package vpn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy"
	"errors"
	"log"
)

func TestHandshake(t *testing.T) {
/*	cs := noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
	staticI := cs.GenerateKeypair(nil)
	staticR := cs.GenerateKeypair(nil)

	h1, err := NewNoiseIKHandshake(
		cs,
		[]byte("caddy-vpn"),
		staticI,
		noise.DHKey{Public:staticR.Public},
		true,
	)
	if err != nil {
		t.Fatal(err)
	}

	h2, err = NewNoiseIKHandshake(
		cs,
		[]byte("caddy-vpn"),
		noise.DHKey{},
		staticR,
		false,
	)*/


}


func TestVPN(t *testing.T) {

	input := `vpn {
		    publickey serverpublickey
		    privatekey serverprivatekey
		    clients {
			publickey client_publickey1
			publickey client_publickey2
			publickey client_publickey3
		    }

		    subnet 192.168.4.1/24
		    mtu 1400
		    dnsport 53
		    auth /auth
		    packet /packet
		}`

	h, err := Parse(caddy.NewTestController("http", input))
	if err != nil {
		t.Fatal(err)
	}

	h.Next = httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
			return 404, errors.New("404 error")
		})


	req, err := http.NewRequest("GET", "http://localhost/auth/", nil)
	if err != nil {
		t.Fatalf("Test: Could not create HTTP request: %v", err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Test fail: code[%d]\n", rec.Code)
	}

	log.Printf("rec.body = %v\n", rec.Body)

}
