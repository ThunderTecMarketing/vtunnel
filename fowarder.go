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
	"github.com/FTwOoO/go-tun2io/tun2io"
	"github.com/FTwOoO/netstack/tcpip"
	"github.com/FTwOoO/netstack/tcpip/link/channel"
	"github.com/FTwOoO/netstack/tcpip/link/sniffer"
	"net"
)

type Fowarder struct {
	linkEP  *channel.Endpoint
	tun2ioM *tun2io.Tun2ioManager
	stack   tcpip.Stack
	//defaultDialer proxy.Dialer
	//dnsServ       *tun2io.DnsServer
}

func NewFowarder(ip net.IP, subnet *net.IPNet) (f *Fowarder, err error) {
	const defaultMTU = 65536
	id, linkEP := channel.New(256, defaultMTU, defaultLinkAddr)
	if false {
		id = sniffer.New(id)
	}

	dialer := new(tun2io.DirectDialer)
	tun2ioM, err := tun2io.Tun2IO(ip, subnet, defaultMTU, true, dialer)
	if err != nil {
		return
	}

	f = &Fowarder{tun2ioM:tun2ioM, stack:tun2ioM.GetStack(), linkEP:linkEP}

	return f, nil

}

func (f *Fowarder) Send([] byte) {

}

func (f *Fowarder) Recv() ([]byte, error) {

}