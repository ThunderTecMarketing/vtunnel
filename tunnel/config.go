/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package tunnel

import (
	"errors"
	"net"
	"io/ioutil"
	"github.com/mholt/caddy"
)

var ConfigFileName = "vtunnel.conf"
var ErrInValidKeyLength = errors.New("Invalid key length")

type TransportType string

var (
	TRANSPORT1 = TransportType("tcp-fragment-gcm256-msgpack")
	handlerGenerators []HandlerGenerator = []HandlerGenerator{}
)

type ListenerHandler func(net.Listener) error
type HandlerGenerator func(config *Config) ListenerHandler

func ResgisterHandlerGenerator(g HandlerGenerator) {
	handlerGenerators = append(handlerGenerators, g)
}

type Config struct {
	ListenAddr     string
	RemoteAddr     string

	IsServer       bool
	TransportKey   []byte
	TransportType  TransportType
	LogFilePath    string
	LogLevel       string

	LocalProxyType string
}

func (s *Config) GetHandler() ListenerHandler {
	for _, g := range handlerGenerators {
		l := g(s)
		if l != nil {
			return l
		}
	}

	return nil
}

func configLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       ConfigFileName,
		ServerTypeName: serverType,
	}, nil
}
