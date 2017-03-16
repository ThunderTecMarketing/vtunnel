package tcpserver

import (
	"net"
	"github.com/mholt/caddy"
	"fmt"
	"strconv"
	"encoding/hex"
)

func SetupTunnelPlugin(c *caddy.Controller) (err error) {

	var m *handler

	if m != nil {
		return c.Err("cannot specify vpn more than once")
	}

	if m, err = Parse(c); err != nil {
		return err
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

	return nil
}

func Parse(c *caddy.Controller) (m *handler, err error) {
	/*
	m = &handler{}

	if c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 1:
			m.VPNPath = args[0]
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
				//c.IncrNest()
				var clientkey []byte

				for c.NextBlock() {
					clientkey, err = HexKey(c.Val(), KeyLength)
					if err != nil {
						break
					}
					m.ClientPublicKeys = append(m.ClientPublicKeys, clientkey)

				}

			default:
				err = c.Errf("Unknown vpn arg: %s", c.Val())
			}
			if err != nil {
				return
			}
		}
	}

	m.DnsServer, err = CreateDnsServer()
	*/
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

	return HexKey(tempkey, keylength)
}

func HexKey(hexkey string, keylength int) (key []byte, err error) {
	key, err = hex.DecodeString(hexkey)
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
