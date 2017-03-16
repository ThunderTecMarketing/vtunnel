package msocks

import (
	"net"
	"github.com/miekg/dns"
)
type Lookuper interface {
	LookupIP(host string) (addrs []net.IP, err error)
}

type NetLookupIP struct {
}

func (n *NetLookupIP) LookupIP(host string) (addrs []net.IP, err error) {
	return net.LookupIP(host)
}

var DefaultLookuper Lookuper

func init() {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return
	}

	var addrs []string
	for _, srv := range conf.Servers {
		addrs = append(addrs, net.JoinHostPort(srv, conf.Port))
	}

	DefaultLookuper = NewDnsLookup(addrs, "")
}

type DnsLookup struct {
	Servers []string
	c       *dns.Client
}

func NewDnsLookup(Servers []string, dnsnet string) (d *DnsLookup) {
	d = &DnsLookup{
		Servers: Servers,
	}
	d.c = new(dns.Client)
	d.c.Net = dnsnet
	return d
}

func (d *DnsLookup) Exchange(m *dns.Msg) (r *dns.Msg, err error) {
	for _, srv := range d.Servers {
		r, _, err = d.c.Exchange(m, srv)
		if err != nil {
			continue
		}
		if len(r.Answer) > 0 {
			return
		}
	}
	return
}

func (d *DnsLookup) query(host string, t uint16, as []net.IP) (addrs []net.IP, err error) {
	addrs = as

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), t)
	m.RecursionDesired = true

	r, err := d.Exchange(m)
	if err != nil {
		return
	}

	for _, a := range r.Answer {
		switch ta := a.(type) {
		case *dns.A:
			addrs = append(addrs, ta.A)
		case *dns.AAAA:
			addrs = append(addrs, ta.AAAA)
		}
	}
	return
}

func (d *DnsLookup) LookupIP(host string) (addrs []net.IP, err error) {
	addrs, err = d.query(host, dns.TypeA, addrs)
	if err != nil {
		return
	}
	addrs, err = d.query(host, dns.TypeAAAA, addrs)
	return
}


// ---- dns part ----

func MakeDnsFrame(host string, t uint16, streamid uint16) (req *dns.Msg, f Frame, err error) {
	log.Debug("make a dns query for %s.", host)

	req = new(dns.Msg)
	req.Id = dns.Id()
	req.SetQuestion(dns.Fqdn(host), t)
	req.RecursionDesired = true

	b, err := req.Pack()
	if err != nil {
		return
	}

	f = &FrameDns{FrameBase.Streamid:streamid, Data:b}
	return
}

func DebugDNS(r *dns.Msg, name string) {
	straddr := ""
	for _, a := range r.Answer {
		switch ta := a.(type) {
		case *dns.A:
			straddr += ta.A.String() + ","
		case *dns.AAAA:
			straddr += ta.AAAA.String() + ","
		}
	}
	log.Info("dns result for %s is %s.", name, straddr)
	return
}

func ParseDnsFrame(f Frame, req *dns.Msg) (addrs []net.IP, err error) {
	ft, ok := f.(*FrameDns)
	if !ok {
		return nil, ErrDnsMsgIllegal
	}

	res := new(dns.Msg)
	err = res.Unpack(ft.Data)
	if err != nil || !res.Response || res.Id != req.Id {
		return nil, ErrDnsMsgIllegal
	}

	if DEBUGDNS {
		DebugDNS(res, req.Question[0].Name)
	}
	for _, a := range res.Answer {
		switch ta := a.(type) {
		case *dns.A:
			addrs = append(addrs, ta.A)
		case *dns.AAAA:
			addrs = append(addrs, ta.AAAA)
		}
	}
	return
}

