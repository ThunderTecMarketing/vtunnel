package caddy-vpn

import (
	"fmt"
	"net/http"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type handler struct {
	Paths []string
	Next  httpserver.Handler
}

func init() {
	caddy.RegisterPlugin("caddy-vpn", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// if the request path is any of the configured paths
	// write hello
	for _, p := range h.Paths {
		if httpserver.Path(r.URL.Path).Matches(p) {
			w.Write([]byte("Hello, I'm a caddy middleware"))
			return 200, nil
		}
	}
	return h.Next.ServeHTTP(w, r)
}

func setup(c *caddy.Controller) error {
	paths, err := parse(c)
	if err != nil {
		return err
	}

	// Runs on Caddy startup, useful for services or other setups.
	c.OnStartup(func() error {
		fmt.Println("caddy-vpn middleware is initiated")
		return nil
	})

	// Runs on Caddy shutdown, useful for cleanups.
	c.OnShutdown(func() error {
		fmt.Println("caddy-vpn middleware is cleaning up")
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
		caddy-vpn /hello
		caddy-vpn /anotherpath
		caddy-vpn {
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

