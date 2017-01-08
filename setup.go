package vpn

import (
	"net"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy"
	"fmt"
	"strconv"
)

func init() {
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
	for c.NextBlock() {
		var err error
		switch c.Val() {
		case "publickey":
			m.PublicKey, err = StringArg(c)
		case "privatekey":
			m.PrivateKey, err = StringArg(c)
		case "clients":
			for c.NextBlock() {
				var pubkey string
				switch c.Val() {
				case "publickey":
					pubkey, err = StringArg(c)
					m.ClientPublicKeys = append(m.ClientPublicKeys, pubkey)
				}
			}
		case "subnet":
			m.Ip, m.Subnet, err = CidrArg(c)
		case "mtu":
			m.MTU, err = UintArg(c)
		case "dnsport":
			m.DnsPort, err = UintArg(c)
		case "auth":
			m.AuthPath, err = StringArg(c)
		case "packet":
			m.PacketPath, err = StringArg(c)
		default:
			return c.Errf("Unknown vpn arg: %s", c.Val())
		}
		if err != nil {
			return err
		}
	}
	return nil
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
		return nil, err
	}
	ip, cidr, err := net.ParseCIDR(a)
	if err != nil {
		return nil, err
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

func UintArg(c *caddy.Controller) (num uint64, err error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return "", c.ArgErr()
	}

	num, err = strconv.ParseUint(args[0], 10, 16)
	return
}
