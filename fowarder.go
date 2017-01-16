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
	"github.com/FTwOoO/netstack/tcpip/buffer"
	"github.com/FTwOoO/netstack/tcpip/network/ipv4"
	"sync"
	"context"
	"errors"
)

type Fowarder struct {
	linkEP      *channel.Endpoint
	tun2ioM     *tun2io.Tun2ioManager
	stack       tcpip.Stack
	//defaultDialer proxy.Dialer
	//dnsServ       *tun2io.DnsServer

	writeChan   chan buffer.View

	readViewsMu sync.Mutex
	readViews   []buffer.View

	ctx         context.Context
	ctxCancel   context.CancelFunc
	closeOne    sync.Once
}

func NewFowarder(ip net.IP, subnet *net.IPNet, mtu int) (f *Fowarder, err error) {
	id, linkEP := channel.New(256, mtu, defaultLinkAddr)
	if false {
		id = sniffer.New(id)
	}

	dialer := new(tun2io.DirectDialer)
	tun2ioM, err := tun2io.Tun2IO(ip, subnet, id, true, dialer)
	if err != nil {
		return
	}

	f = &Fowarder{
		tun2ioM:tun2ioM,
		stack:tun2ioM.GetStack(),
		linkEP:linkEP,
		writeChan:make(chan []byte, 1024),
	}

	f.ctx, f.ctxCancel = context.WithCancel(context.Background())
	go f.reader()
	go f.writer()

	return f, nil
}

func (f *Fowarder) writer() {
	for {
		select {
		case b := <-f.writeChan:
			views := [1]buffer.View{b}
			vv := buffer.NewVectorisedView(len(b), views)
			f.linkEP.Inject(ipv4.ProtocolNumber, &vv)
		case <-f.ctx.Done():
			return
		}
	}
}

func (f *Fowarder) reader() {

	for {
		select {
		case b := <-f.linkEP.C:
			f.readViewsMu.Lock()
			f.readViews = append(f.readViews, buffer.View(b))
			f.readViewsMu.Unlock()
		case <-f.ctx.Done():
			return
		}
	}

}

func (f *Fowarder) Send(b buffer.View) {
	f.writeChan <- buffer.View(b)
}

func (f *Fowarder) Recv() ([]buffer.View, error) {

	f.readViewsMu.Lock()
	defer f.readViewsMu.Unlock()

	if len(f.readViews) > 0 {
		ret := f.readViews
		f.readViews = []buffer.View{}
		return ret, nil
	}

	return nil, errors.New("???")
}


func(f *Fowarder) Close(reason error) {

	f.closeOne.Do(func() {
		f.ctxCancel()
	})
	return
}