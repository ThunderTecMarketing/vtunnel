package msocks

import (
	"errors"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/dnsrelay/dnsrelay"
)

type MsocksServer struct {
	*SessionPool
	dialer    Dialer
	dnsServer *dnsrelay.DNSServer
}

func NewMsocksServer(dialer Dialer) (ms *MsocksServer, err error) {
	if dialer == nil {
		err = errors.New("empty dialer")
		log.Errorf("%s", err)
		return
	}

	dnsServer, err := dnsrelay.NewDNSServer(nil, true)
	if err != nil {
		return
	}

	ms = &MsocksServer{
		SessionPool: CreateSessionPool(0, 0),
		dialer:      dialer,
		dnsServer: dnsServer,
	}

	return
}

func (ms *MsocksServer) Handler(conn conn.ObjectIO) {
	//log.Notice("connection come from: %s => %s.", conn.RemoteAddr(), conn.LocalAddr())

	sess := NewSession(conn, ms.dnsServer)
	sess.next_id = 1
	sess.dialer = ms.dialer

	ms.Add(sess)
	defer ms.Remove(sess)
	sess.Run()

	//log.Notice("server session %d quit: %s => %s.", sess.LocalPort(), conn.RemoteAddr(), conn.LocalAddr())
}

func (ms *MsocksServer) Serve(listener conn.ObjectListener) (err error) {
	var conn conn.ObjectIO

	for {
		conn, err = listener.Accept()
		if err != nil {
			log.Errorf("%s", err)
			continue
		}
		go func() {
			defer conn.Close()
			ms.Handler(conn)
		}()
	}
	return
}
