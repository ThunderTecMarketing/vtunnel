/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package tunnel

import (
	"time"
	"errors"
	"net"
)

var DefaultPrologue = "vtunnel"
var KeyLength = 16

var DefaultPeerTimeout = time.Duration(30 * time.Second)
var DefaultTokenTimeout = time.Duration(10 * time.Second)

var ErrInValidHandshakeStep = errors.New("Invalid handshake step")
var ErrInValidKeyLength = errors.New("Invalid key length")
var ErrPeerAlreadyExist = errors.New("Peer exists")
var ErrPacketLengthInvalid = errors.New("Packet length is not in (0,MTU]")
var ErrWriteFail = errors.New("Write fail")

type TransportType string

var (
	TRANSPORT1 = TransportType("TCP-Fragment-AheadGCM256-Msgpack")
	handlerGenerators []HandlerGenerator = []HandlerGenerator{}
)

type ListenerHandler func(net.Listener) error
type HandlerGenerator func (config *Config) ListenerHandler

func ResgisterHandlerGenerator(g HandlerGenerator) {
	handlerGenerators = append(handlerGenerators, g)
}

type Config struct {
	ListenAddr     string
	RemoteAddr     string

	IsServer       bool
	TransportKey   []byte
	TransportType  TransportType
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
