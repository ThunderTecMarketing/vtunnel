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

	var tempkey string
	var clientkey []byte

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
				tempkey, err = StringArg(c)
				if err == nil {
					m.PublicKey, err = hex.DecodeString(tempkey)
				}

			case "privatekey":
				tempkey, err = StringArg(c)
				if err == nil {
					m.PrivateKey, err = hex.DecodeString(tempkey)
				}
			case "clients":
				c.Next()
				c.IncrNest()
				for c.NextBlock() {
					switch c.Val() {
					case "publickey":
						tempkey, err = StringArg(c)
						if err != nil {
							return
						}

						clientkey, err =  hex.DecodeString(tempkey)
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

	m.Fowarder, err = NewFowarder(m.Ip, m.Subnet, m.MTU)
	return
}


// Helpers below here could potentially be methods on *caddy.Contoller for convenience

// Assert only one arg and return it
func StringArg(c *caddy.Controller) (string, error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return "", c.ArgErr()
	}
	return args[0], nil
}

// Assert only one arg is a valid cidr notation
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
