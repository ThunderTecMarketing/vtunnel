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
	"fmt"
	"encoding/hex"
	"net/url"
	"strconv"
)

var DefaultLinkAddr = tcpip.LinkAddress("\x0a\x0a\x0b\x0b\x0c\x0c")
var DefaultPrologue = "caddy-vpn"
var DefaultCipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
var KeyLength = noise.DH25519.DHLen()

var DefaultPeerTimeout = time.Duration(30 * time.Second)
var DefaultTokenTimeout = time.Duration(10 * time.Second)

var ErrInValidHandshakeStep = errors.New("Invalid handshake step")
var ErrInValidKeyLength = errors.New("Invalid key length")
var ErrPeerAlreadyExist = errors.New("Peer exists")
var ErrPacketLengthInvalid = errors.New("Packet length is not in (0,MTU]")

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

type ClientSetting struct {
	Ip     net.IP
	Subnet *net.IPNet
	Mtu    uint16
	Token  []byte
}

func (cs *ClientSetting) Encode() string {
	maskNum, _ := cs.Subnet.Mask.Size()

	clientSetting := make(url.Values)
	clientSetting.Set("ip", fmt.Sprintf("%s/%d", cs.Ip.String(), maskNum))
	clientSetting.Set("mtu", fmt.Sprintf("%d", cs.Mtu))
	clientSetting.Set("token", hex.EncodeToString(cs.Token))

	return clientSetting.Encode()
}

func DecodeClientSetting(content string) (cs *ClientSetting, err error) {
	q, err := url.ParseQuery(content)
	if err != nil {
		return
	}

	cs = new(ClientSetting)
	ipandnet := q.Get("ip")
	if ipandnet == "" {
		return nil, errors.New("Need ip")
	}

	cs.Ip, cs.Subnet, err = net.ParseCIDR(ipandnet)
	if err != nil {
		return
	}

	mtuS := q.Get("mtu")
	if ipandnet == "" {
		return nil, errors.New("Need mtu")
	}

	mtu, err := strconv.ParseInt(mtuS, 10, 16)
	cs.Mtu = uint16(mtu)

	tokenS := q.Get("token")
	if tokenS == "" {
		return nil, errors.New("Need token")
	}

	token, err := hex.DecodeString(tokenS)
	if err != nil {
		return
	}
	cs.Token = token

	return
}