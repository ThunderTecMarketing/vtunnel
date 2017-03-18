/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package tcpserver

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


type ServerListenerHandler func (net.Listener)

type ServerConfig struct {
	ListenHost string
	ListenPort uint16

	Clients map[string]string
	Handler ServerListenerHandler

}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{Clients: make(map[string]string)}
}