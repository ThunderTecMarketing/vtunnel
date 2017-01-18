package vpn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy"
	"github.com/FTwOoO/noise"
	"errors"
	"log"
	"encoding/hex"
	"bytes"
)

func getClientHandshake(ServerPublicKey []byte) (h *NoiseIKHandshake, err error) {
	cs := noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)

	publicKey, _:= hex.DecodeString("04537cd141acdc2feba13b623b2c3f6151cad48384fd6cc8065399dcdd2d257d")
	privateKey,_:= hex.DecodeString("c0f2adf5c07b865b9b615eebafc352954ac4dd7b0d4bd55499880e3b7fd05448")
	staticI := noise.DHKey{Public:publicKey, Private:privateKey}

	h, err = NewNoiseIKHandshake(
		cs,
		[]byte(DefaultPrologue),
		staticI,
		noise.DHKey{Public:ServerPublicKey},
		true,
	)
	return
}


func TestHandshake(t *testing.T) {
	serverPublicKey, _:= hex.DecodeString("e8e394b473b7b58514404fdddc0dd237ff631ceba3c0d1eddcddecb58f5a7d2a")
	clientHandshake, err := getClientHandshake(serverPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	input := `vpn {
		    publickey  e8e394b473b7b58514404fdddc0dd237ff631ceba3c0d1eddcddecb58f5a7d2a
		    privatekey 3fbf4c6e081f845ab7998471dd4af084eea403f66a87cb5c2d775fbaa6c76eb4
		    clients {
			publickey 04537cd141acdc2feba13b623b2c3f6151cad48384fd6cc8065399dcdd2d257d
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


	reqContent := "test"
	encodedReqContent, err := clientHandshake.Encode([]byte(reqContent))

	req, err := http.NewRequest("GET", "http://localhost/auth/", bytes.NewBuffer(encodedReqContent))
	if err != nil {
		t.Fatalf("Test: Could not create HTTP request: %v", err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Test fail: code[%d]\n", rec.Code)
	}

	respContent := make([]byte, 1024)
	n, err := rec.Body.Read(respContent)
	decodedRespContent, err := clientHandshake.Decode(respContent[:n])

	log.Printf("Get response = %v\n", decodedRespContent)

}
