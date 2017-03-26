/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package tunnel

import (
	"errors"
	"net"
)

var ErrInValidKeyLength = errors.New("Invalid key length")

type TransportType string

var (
	TRANSPORT1 = TransportType("TCP-Fragment-AheadGCM256-Msgpack")
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
	LogLevel   string

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
