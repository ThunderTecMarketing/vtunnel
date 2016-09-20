package vpn

import (
	"fmt"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)



func init() {
	fmt.Println("init!!!!!!.....")
	caddy.RegisterPlugin("vpn", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	fmt.Println("setup!!!!!!.....")

	paths, err := parse(c)
	if err != nil {
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

	config := httpserver.GetConfig(c)

	config.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return &handler{
			Paths: paths,
			Next:  next,
		}
	})

	return nil
}

func parse(c *caddy.Controller) ([]string, error) {
	// This parses the following config blocks
	/*
		vpn /hello
		vpn /anotherpath
		vpn {
			path /hello
			path /anotherpath
		}
	*/
	var paths []string
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			// no argument passed, check the config block
			for c.NextBlock() {
				switch c.Val() {
				case "path":
					if !c.NextArg() {
						// we are expecting a value
						return paths, c.ArgErr()
					}
					p := c.Val()
					paths = append(paths, p)
					if c.NextArg() {
						// we are expecting only one value.
						return paths, c.ArgErr()
					}
				}
			}
		case 1:
			// one argument passed
			paths = append(paths, args[0])
			if c.NextBlock() {
				// path specified, no block required.
				return paths, c.ArgErr()
			}
		default:
			// we want only one argument max
			return paths, c.ArgErr()
		}
	}
	return paths, nil
}

