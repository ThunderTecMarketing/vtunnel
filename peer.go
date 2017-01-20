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
	"time"
	"sync"
	"github.com/FTwOoO/vpncore/net/tcpip"
)

type Peer struct {
	PublicKey        []byte
	token            string
	Ip               net.IP
	NoiseIKHandshake *NoiseIXHandshake

	HandshakeDone    bool
	lastSeenTime     time.Time
}

func (peer *Peer) Touch() {
	peer.lastSeenTime = time.Now()
}

type Peers struct {
	config      *Config
	MyIp        net.IP
	IpPool      *tcpip.IP4Pool

	PeerTimeout chan *Peer

	peerByIp    map[string]*Peer
	peerByKey   map[string]*Peer
	peerLock    sync.RWMutex
}

func NewPeers(config *Config) (vs *Peers) {
	vs = new(Peers)
	vs.config = config
	vs.IpPool, _ = tcpip.NewIP4Pool(config.Subnet)
	vs.MyIp = config.Ip

	vs.peerByIp = map[string]*Peer{}
	vs.peerByKey = map[string]*Peer{}
	vs.PeerTimeout = make(chan *Peer, 100)

	go vs.checkTimeout(DefaultPeerTimeout)
	return
}

func (vs *Peers) checkTimeout(timeout time.Duration) {
	vs.peerLock.RLock()
	defer vs.peerLock.RUnlock()

	for _, peer := range vs.peerByIp {
		conntime := time.Since(peer.lastSeenTime)
		if conntime > timeout {
			vs.PeerTimeout <- peer
		}
	}
}

func (vs *Peers) AddPeer(publicKey []byte, handshake *NoiseIXHandshake) (peer *Peer, err error) {
	vs.peerLock.RLock()
	defer vs.peerLock.RUnlock()

	if _, ok := vs.peerByKey[string(publicKey)]; ok {
		return nil, ErrPeerAlreadyExist
	}

	var ip net.IP
	for {
		ip, err = vs.IpPool.Next()
		if err != nil {
			return
		}

		if ip.Equal(vs.MyIp) {
			continue
		}

		break
	}

	peer = &Peer{Ip:ip, PublicKey:publicKey, NoiseIKHandshake: handshake}

	vs.peerByIp[peer.Ip.String()] = peer
	vs.peerByKey[string(publicKey)] = peer
	return
}

func (vs *Peers) DeletePeer(peer *Peer) {
	vs.peerLock.RLock()
	defer vs.peerLock.RUnlock()

	vs.IpPool.Release(peer.Ip)
	delete(vs.peerByIp, peer.Ip.String())
	delete(vs.peerByKey, string(peer.PublicKey))

}

func (vs *Peers) GetPeerByIp(ip net.IP) (*Peer) {
	vs.peerLock.RLock()
	defer vs.peerLock.RUnlock()

	return vs.peerByIp[ip.String()]
}

func (vs *Peers) GetPeerById(id []byte) (*Peer) {
	vs.peerLock.RLock()
	defer vs.peerLock.RUnlock()

	return vs.peerByKey[string(id)]
}