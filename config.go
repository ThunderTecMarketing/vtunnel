/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package vpn

import (
	"net"
	"github.com/FTwOoO/netstack/tcpip"
	"errors"
	"github.com/FTwOoO/noise"
	"time"
)

var DefaultLinkAddr = tcpip.LinkAddress("\x0a\x0a\x0b\x0b\x0c\x0c")
var DefaultPrologue = "caddy-vpn"
var DefaultCipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
var KeyLength = noise.DH25519.DHLen()

var DefaultPeerTimeout = time.Duration(30 *time.Second)
var DefaultTokenTimeout = time.Duration(10*time.Second)

var ErrInValidHandshakeStep = errors.New("Invalid handshake step")
var ErrInValidKeyLength = errors.New("Invalid key length")
var ErrPeerAlreadyExist = errors.New("peer exists")

type Config struct {
	PublicKey        []byte
	PrivateKey       []byte
	ClientPublicKeys [][]byte
	Ip               net.IP
	Subnet           *net.IPNet
	MTU              uint16
	DnsPort          uint16
	AuthPath         string
	PacketPath       string
}
