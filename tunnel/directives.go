package tunnel

import (
	"github.com/mholt/caddy"
)

var directives = []string{
	"client",
	"server",
}

func init() {
	caddy.RegisterPlugin("client", caddy.Plugin{
		ServerType: ServerType,
		Action:     SetupDirective,
	})

	caddy.RegisterPlugin("server", caddy.Plugin{
		ServerType: ServerType,
		Action:     SetupDirective,
	})
}

func SetupDirective(c *caddy.Controller) (err error) {

	ctx := c.Context().(*tunnelContext)
	config := ctx.keysToConfigs[c.Key]

	if c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			break
		default:
			return c.ArgErr()
		}

		switch c.Val() {
		case "server", "client":
			if c.Val() == "client" {
				config.IsServer = false
			} else {
				config.IsServer = true
			}

			for c.NextBlock() {
				switch c.Val() {
				case "proxyType":
					config.LocalProxyType, err = StringArg(c)
				case "remoteAddr":
					config.RemoteAddr, err = StringArg(c)

				case "transportType":
					var transportType string
					transportType, err = StringArg(c)
					config.TransportType = TransportType(transportType)
				case "transportKey":
					var transportKey string
					transportKey, err = StringArg(c)
					config.TransportKey = []byte(transportKey)

				}
			}

		default:
			break
		}

	}

	return
}
