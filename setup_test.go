package vpn

import (
	"reflect"
	"testing"
	"github.com/mholt/caddy"
	"net"
)

func TestHeadersParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		expected  *handler
	}{

		{`realip {
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
		}`},
		false,
		&handler{
			PublicKey:"serverpublickey",
			PrivateKey:"serverprivatekey",
			ClientPublicKeys: []string{"client_publickey1", "client_publickey2", "client_publickey3"},
			Ip:net.IP{192, 168, 4, 1},
			Subnet:net.IPNet{IP:net.IP{192, 168, 4, 0}, Mask: net.CIDRMask(24, 8)},
			MTU:1400,
			DnsPort:53,
			AuthPath:"/auth",
			PacketPath:"/packet",
		},
	}

	for i, test := range tests {
		actual, err := Parse(caddy.NewTestController("http", test.input))

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if len(actual) != len(test.expected) {
			t.Fatalf("Test %d expected %d rules, but got %d", i, len(test.expected), len(actual))
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, but got %v", test.expected, actual)
		}

	}
}
