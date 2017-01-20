package vpn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy"
	"github.com/FTwOoO/noise"
	"errors"
	"encoding/hex"
	"bytes"
	"fmt"
)

var validClientPublicKey = "e01ee3207ea15d346c362b7e20cef3a1088ec0a11a1141b3584ed44e2bb69531"
var validClientPrivateKey = "22e70850eb2da8fe184ed4998575f403f24c7ad54dbdc2132ae6a44c81b41180"

var invalidClientPublicKey = "5561cbf77dc96a21041f2a6127b927e439a15852dd4a915e7741fe4889afdb34"
var invalidClientPrivateKey = "d15fde7c16da6364374e8cc96f934c13de5b686aa0ad005e5ba2093fe2ff5da3"

var serverPublicKey = "e8e394b473b7b58514404fdddc0dd237ff631ceba3c0d1eddcddecb58f5a7d2a"
var serverPrivateKey = "3fbf4c6e081f845ab7998471dd4af084eea403f66a87cb5c2d775fbaa6c76eb4"


func getClientHandshake(pub string, pri string) (h *NoiseIXHandshake, err error) {
	cs := noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
	publicKey, _ := hex.DecodeString(pub)
	privateKey, _ := hex.DecodeString(pri)
	staticI := noise.DHKey{Public:publicKey, Private:privateKey}

	h, err = NewNoiseIXHandshake(
		cs,
		[]byte(DefaultPrologue),
		staticI,
		true,
	)
	return
}

func TestHandshake(t *testing.T) {

	input := fmt.Sprintf(`vpn {
		    publickey %s
		    privatekey %s
		    clients {
			publickey %s
		    }

		    subnet 192.168.4.1/24
		    mtu 1400
		    dnsport 53
		    auth /auth
		    packet /packet
		}`, serverPublicKey, serverPrivateKey, validClientPublicKey)

	h, err := Parse(caddy.NewTestController("http", input))
	if err != nil {
		t.Fatal(err)
	}

	h.Next = httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusNotFound, errors.New("404 error")
	})

	validH, err := getClientHandshake(validClientPublicKey, validClientPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	invalidH, err := getClientHandshake(invalidClientPublicKey, invalidClientPrivateKey)
	if err != nil {
		t.Fatal(err)
	}


	Handshake(t, h, validH, http.StatusOK)
	Handshake(t, h, invalidH, http.StatusUnauthorized)

}

func Handshake(t *testing.T, h *handler, clientHandshake *NoiseIXHandshake, expectedCode int) {
	reqContent := []byte("test")
	encodedReqContent, err := clientHandshake.Encode(reqContent)

	req, err := http.NewRequest("GET", "http://localhost/auth/", bytes.NewBuffer(encodedReqContent))
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	rec := httptest.NewRecorder()
	statusCode, err := h.ServeHTTP(rec, req)
	if statusCode != expectedCode {
		t.Fatalf("Code [%d] not expected[%d]\n", statusCode, expectedCode)
	}

	if expectedCode == http.StatusOK {
		respContent := make([]byte, 1024)
		n, err := rec.Body.Read(respContent)
		if err != nil {
			t.Fatal(err)
		}

		_, err = clientHandshake.Decode(respContent[:n])
		if err != nil {
			t.Fatal(err)
		}


	}
}