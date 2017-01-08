package vpn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy"
)

func TestRealIP(t *testing.T) {

	input := `realip {
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
			return 0, nil
		})


	req, err := http.NewRequest("GET", "http://ftwo.me/auth/", nil)
	if err != nil {
		t.Fatalf("Test: Could not create HTTP request: %v", err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Test fail: code[%d]\n", rec.Code)
	}

}
