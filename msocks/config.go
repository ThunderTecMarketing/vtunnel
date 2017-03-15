package msocks

import (
	"errors"
	"math/rand"
	"time"

	logging "github.com/op/go-logging"
	"net"
	"github.com/miekg/dns"
)

const (
	DIAL_RETRY   = 2
	DIAL_TIMEOUT = 30
	DNS_TIMEOUT  = 30

	SHRINK_TIME = 3
	DEBUGDNS    = false
)

const (
	ERR_NONE = iota
	ERR_AUTH
	ERR_IDEXIST
	ERR_CONNFAILED
	ERR_TIMEOUT
	ERR_CLOSED
)

var (
	ErrNoSession       = errors.New("session in pool but can't pick one.")
	ErrSessionNotFound = errors.New("session not found.")
	ErrAuthFailed      = errors.New("auth failed.")
	ErrAuthTimeout     = errors.New("auth timeout %s.")
	ErrStreamNotExist  = errors.New("stream not exist.")
	ErrQueueClosed     = errors.New("queue closed.")
	ErrUnexpectedPkg   = errors.New("unexpected package.")
	ErrNotSyn          = errors.New("frame result in conn which status is not syn.")
	ErrFinState        = errors.New("status not est or fin wait when get fin.")
	ErrIdExist         = errors.New("frame sync stream id exist.")
	ErrState           = errors.New("status error.")
	ErrUnknownState    = errors.New("unknown status.")
	ErrChanClosed      = errors.New("chan closed.")
	ErrDnsTimeOut      = errors.New("dns timeout.")
	ErrDnsMsgIllegal   = errors.New("dns message illegal.")
	ErrNoDnsServer     = errors.New("no proper dns server.")
)

var (
	log = logging.MustGetLogger("msocks")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}


type Dialer interface {
	Dial(string, string) (net.Conn, error)
}


type TcpDialer struct {
}

func (td *TcpDialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

func (td *TcpDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

var DefaultTcpDialer = &TcpDialer{}

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