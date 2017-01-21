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
	"github.com/FTwOoO/netstack/tcpip/link/sniffer"
	"github.com/FTwOoO/netstack/tcpip/buffer"
	"net"
	"sync"
	"context"
	"github.com/FTwOoO/netstack/tcpip/network/ipv4"
	"github.com/FTwOoO/netstack/tcpip/header"
	"log"
	"github.com/FTwOoO/netstack/tcpip/link/channel"
)

type Fowarder struct {
	ip          net.IP
	subnet      *net.IPNet
	linkEP      *channel.Endpoint
	tun2ioM     *tun2io.Tun2ioManager
	stack       tcpip.Stack
	//defaultDialer proxy.Dialer
	//dnsServ       *tun2io.DnsServer

	writeChan   chan buffer.View

	readViewsMu sync.Mutex
	readViews   map[string][]buffer.View

	ctx         context.Context
	ctxCancel   context.CancelFunc
	closeOne    sync.Once
}

func NewFowarder(ip net.IP, subnet *net.IPNet, mtu uint16) (f *Fowarder, err error) {
	id, linkEP := channel.New(256, uint32(mtu), DefaultLinkAddr)
	if false {
		id = sniffer.New(id)
	}

	dialer := new(tun2io.DirectDialer)
	tun2ioM, err := tun2io.Tun2IO(ip, subnet, id, true, dialer)
	if err != nil {
		return
	}

	f = &Fowarder{
		ip:ip,
		subnet:subnet,
		tun2ioM:tun2ioM,
		stack:tun2ioM.GetStack(),
		linkEP:linkEP,
		writeChan:make(chan buffer.View, 1024),
		readViews:make(map[string][]buffer.View),
	}

	f.ctx, f.ctxCancel = context.WithCancel(context.Background())
	go f.reader()
	go f.writer()

	return f, nil
}

func (f *Fowarder) writer() {
	ipv4Proto := ipv4.NewProtocol()
	//ipv6Proto := ipv6.NewProtocol()

	Writing:
	for {
		select {
		case b := <-f.writeChan:

			switch header.IPVersion(b) {
			case header.IPv4Version:
				src, dst := ipv4Proto.ParseAddresses(b)
				srcIP := net.ParseIP(src.String())
				dstIp := net.ParseIP(dst.String())

				if !f.subnet.Contains(srcIP) {
					log.Printf("Ip src not allowed: %s", srcIP.String())
					continue Writing
				}

				if f.subnet.Contains(dstIp) && !dstIp.Equal(f.ip) {
					f.pushPacketToTarget(b, dstIp)
					continue Writing
				}

				views := []buffer.View{b}
				vv := buffer.NewVectorisedView(len(b), views)

				f.linkEP.Inject(header.IPv4ProtocolNumber, &vv)

			case header.IPv6Version:
				//TODO: support ipv6...
				continue Writing
			default:
				log.Printf("Bad ip packet %x\n", b)
				continue Writing
			}



		case <-f.ctx.Done():
			return
		}
	}
}

func (f *Fowarder) pushPacketToTarget(b []byte, dst net.IP) {
	f.readViewsMu.Lock()
	key := dst.String()

	if _, ok := f.readViews[key]; !ok {
		f.readViews[key] = []buffer.View{}
	}

	f.readViews[key] = append(f.readViews[key], b)
	f.readViewsMu.Unlock()
}

func (f *Fowarder) reader() {

	for {
		select {
		case p := <-f.linkEP.C:
			newPacket := append([]byte(p.Header), []byte(p.Payload)...)

			targetIp := net.ParseIP(p.Route.RemoteAddress)
			f.pushPacketToTarget(newPacket, targetIp)

		case <-f.ctx.Done():
			return
		}
	}

}

func (f *Fowarder) Send(b buffer.View) {
	f.writeChan <- buffer.View(b)
}

func (f *Fowarder) Recv(dst net.IP) ([]buffer.View) {

	f.readViewsMu.Lock()
	defer f.readViewsMu.Unlock()

	ip := dst.String()

	if _, ok := f.readViews[ip]; !ok || len(f.readViews) <= 0 {
		return nil
	}

	if len(f.readViews) > 0 {
		ret := f.readViews[ip]
		f.readViews = []buffer.View{}
		return ret
	}

	return nil
}

func (f *Fowarder) DeleteTarget(dst net.IP) {
	f.readViewsMu.Lock()
	defer f.readViewsMu.Unlock()

	delete(f.readViews, dst.String())
}

func (f *Fowarder) Close(reason error) {

	f.closeOne.Do(func() {
		f.ctxCancel()
	})
	return
}