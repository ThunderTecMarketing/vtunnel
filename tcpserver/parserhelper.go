package tcpserver

import (
	"github.com/mholt/caddy"
	"strconv"
	"encoding/hex"
)

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

