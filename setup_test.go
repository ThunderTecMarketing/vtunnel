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

		{
			`vpn {
			    publickey 0001020304
			    privatekey  0001020304
			    clients {
				publickey 0001020304
				publickey 0001020304
				publickey 0001020304
			    }

			    subnet 192.168.4.1/24
			    mtu 1400
			    dnsport 53
			    auth /auth
			    packet /packet
			}`,

			false,
			&handler{
				Config:Config{
					PublicKey:[]byte{
						0x00, 0x01, 0x02, 0x03, 0x04,
					},
					PrivateKey:[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
					ClientPublicKeys: [][]byte{
						[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
						[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
						[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
					},
					Ip:net.IPv4(192, 168, 4, 1).To4(),
					Subnet:&net.IPNet{IP:net.IPv4(192, 168, 4, 0).To4(), Mask: net.CIDRMask(24, 32)},
					MTU:1400,
					DnsPort:53,
					AuthPath:"/auth",
					PacketPath:"/packet",
				},
			},
		},
	}

	for i, test := range tests {
		actual, err := Parse(caddy.NewTestController("http", test.input))

		if err == nil && test.shouldErr {
			t.Fatalf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Fatalf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if !reflect.DeepEqual(test.expected.Config, actual.Config) {
			t.Fatalf("Expected %v, but got %v", test.expected, actual)
		}

	}
}
