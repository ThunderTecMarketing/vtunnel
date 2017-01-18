package vpn

import (
	"net"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy"
	"fmt"
	"strconv"
	"log"
	"os"
	"encoding/hex"
	"github.com/FTwOoO/noise"
)

func init() {

	if os.Getenv("CADDY_DEV_MODE") == "1" {
		httpserver.RegisterDevDirective("vpn", "")
	}
	caddy.RegisterPlugin("vpn", caddy.Plugin{
		ServerType: "http",
		Action:     Setup,
	})

}

func Setup(c *caddy.Controller) (err error) {
	var m *handler

	for c.Next() {
		if m != nil {
			return c.Err("cannot specify vpn more than once")
		}

		if m, err = Parse(c); err != nil {
			return err
		}
	}

	// Runs on Caddy startup, useful for services or other setups.
	c.OnStartup(func() error {
		fmt.Println("vpn middleware is initiated")
		return nil
	})

	// Runs on Caddy shutdown, useful for cleanups.
	c.OnShutdown(func() error {
		fmt.Println("vpn middleware is cleaning up")
		return nil
	})

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		m.Next = next
		return m
	})
	return nil
}

func Parse(c *caddy.Controller) (m *handler, err error) {
	m = &handler{}

	if c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			break
		default:
			return nil, c.ArgErr()
		}

		for c.NextBlock() {
			switch c.Val() {
			case "publickey":
				m.PublicKey, err = HexKeyArg(c, KeyLength)
			case "privatekey":
				m.PrivateKey, err = HexKeyArg(c, KeyLength)
			case "clients":
				c.Next()
				c.IncrNest()
				var clientkey []byte

				for c.NextBlock() {
					switch c.Val() {
					case "publickey":
						clientkey, err = HexKeyArg(c, KeyLength)
						m.ClientPublicKeys = append(m.ClientPublicKeys, clientkey)
					default:
						log.Print("Error publickey now\n")
						return nil, c.ArgErr()
					}
				}

			case "subnet":
				m.Ip, m.Subnet, err = CidrArg(c)
				if m.Ip.To4() != nil {
					m.Ip = m.Ip.To4()
				}

				if m.Subnet.IP.To4() != nil {
					m.Subnet.IP = m.Subnet.IP.To4()
				}
			case "mtu":
				m.MTU, err = Uint16Arg(c)
			case "dnsport":
				m.DnsPort, err = Uint16Arg(c)
			case "auth":
				m.AuthPath, err = StringArg(c)
			case "packet":
				m.PacketPath, err = StringArg(c)
			default:
				err = c.Errf("Unknown vpn arg: %s", c.Val())
			}
			if err != nil {
				return
			}
		}
	}

	cipherSuite := DefaultCipherSuite
	m.NoiseIKHandshake, err = NewNoiseIKHandshake(
		cipherSuite,
		[]byte(DefaultPrologue),
		noise.DHKey{},
		noise.DHKey{Public:m.PublicKey, Private:m.PrivateKey},
		false,
	)
	m.Fowarder, err = NewFowarder(m.Ip, m.Subnet, m.MTU)
	return
}

func StringArg(c *caddy.Controller) (string, error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return "", c.ArgErr()
	}
	return args[0], nil
}

func HexKeyArg(c *caddy.Controller, keylength int) (key []byte, err error) {
	tempkey, err := StringArg(c)
	if err != nil {
		return
	}

	key, err = hex.DecodeString(tempkey)
	if err != nil {
		return
	}

	if len(key) != keylength {
		err = ErrInValidKeyLength
		return
	}

	return
}

func CidrArg(c *caddy.Controller) (net.IP, *net.IPNet, error) {
	a, err := StringArg(c)
	if err != nil {
		return []byte{}, nil, err
	}
	ip, cidr, err := net.ParseCIDR(a)
	if err != nil {
		return []byte{}, nil, err
	}
	return ip, cidr, nil
}

func BoolArg(c *caddy.Controller) (bool, error) {
	args := c.RemainingArgs()
	if len(args) > 1 {
		return false, c.ArgErr()
	}
	if len(args) == 0 {
		return true, nil
	}
	switch args[0] {
	case "false":
		return false, nil
	case "true":
		return true, nil
	default:
		return false, c.Errf("Unexpected bool value: %s", args[0])
	}
}

func NoArgs(c *caddy.Controller) error {
	if len(c.RemainingArgs()) != 0 {
		return c.ArgErr()
	}
	return nil
}

func Uint16Arg(c *caddy.Controller) (num uint16, err error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return 0, c.ArgErr()
	}

	num64, err := strconv.ParseUint(args[0], 10, 16)
	num = uint16(num64)
	return
}
