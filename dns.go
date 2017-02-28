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
	"github.com/miekg/dns"
	"net"
	"github.com/FTwOoO/dnsrelay/dnsrelay"
)

type writeFunc func (m []byte) (int, error)

type sessionWriter struct {
	writeFunc writeFunc
}

// WriteMsg implements the ResponseWriter.WriteMsg method.
func (w *sessionWriter) WriteMsg(m *dns.Msg) (err error) {
	var data []byte
	data, err = m.Pack()
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// Write implements the ResponseWriter.Write method.
func (w *sessionWriter) Write(m []byte) (int, error) {
	return w.writeFunc(m)

}

// LocalAddr implements the ResponseWriter.LocalAddr method.
func (w *sessionWriter) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr implements the ResponseWriter.RemoteAddr method.
func (w *sessionWriter) RemoteAddr() net.Addr {
	return nil
}

// TsigStatus implements the ResponseWriter.TsigStatus method.
func (w *sessionWriter) TsigStatus() error {
	return nil
}

// TsigTimersOnly implements the ResponseWriter.TsigTimersOnly method.
func (w *sessionWriter) TsigTimersOnly(b bool) {}

// Hijack implements the ResponseWriter.Hijack method.
func (w *sessionWriter) Hijack() {}

// Close implements the ResponseWriter.Close method
func (w *sessionWriter) Close() error {
	return nil
}

type DnsServer struct {
	h dns.Handler
}

func CreateDnsServer() (*DnsServer, error) {
	handlerServ, err := dnsrelay.NewDNSServer(nil, true)
	if err != nil {
		return nil, err
	}

	d := &DnsServer{h:handlerServ}

	return d, nil
}

func (d *DnsServer) query(data []byte, writeFunc writeFunc) {

	w := &sessionWriter{writeFunc:writeFunc}

	req := new(dns.Msg)
	err := req.Unpack(data)
	if err != nil {
		x := new(dns.Msg)
		x.SetRcodeFormatError(req)
		w.WriteMsg(x)
		w.Close()
	}

	d.h.ServeDNS(w, req)
}