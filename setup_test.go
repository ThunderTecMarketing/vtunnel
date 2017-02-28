package vpn

import (
	"reflect"
	"testing"
	"github.com/mholt/caddy"
	"encoding/hex"
	"fmt"
)

func TestHeadersParse(t *testing.T) {

	defaultHexKey := "000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"
	defaultKey, _ := hex.DecodeString(defaultHexKey)

	tests := []struct {
		input     string
		shouldErr bool
		expected  *handler
	}{

		{
			fmt.Sprintf(`vpn /vpn {
			    publickey %s
			    privatekey  %s
			    clients {
				%s
				%s
				%s
			    }

			}`, defaultHexKey, defaultHexKey, defaultHexKey, defaultHexKey, defaultHexKey),

			false,
			&handler{
				Config:Config{
					PublicKey:defaultKey,
					PrivateKey:defaultKey,
					ClientPublicKeys: [][]byte{
						defaultKey,
						defaultKey,
						defaultKey,
					},

					VPNPath:"/vpn",
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
