package msocks

import (
	"errors"
	"net"
)

type MsocksServer struct {
	*SessionPool
	dialer   Dialer
}

func NewServer(dialer Dialer) (ms *MsocksServer, err error) {
	if dialer == nil {
		err = errors.New("empty dialer")
		log.Error("%s", err)
		return
	}
	ms = &MsocksServer{
		SessionPool: CreateSessionPool(0, 0),
		dialer:      dialer,
	}

	return
}

func (ms *MsocksServer) Handler(conn net.Conn) {
	log.Notice("connection come from: %s => %s.", conn.RemoteAddr(), conn.LocalAddr())

	sess := NewSession(conn)
	sess.next_id = 1
	sess.dialer = ms.dialer

	ms.Add(sess)
	defer ms.Remove(sess)
	sess.Run()

	log.Notice("server session %d quit: %s => %s.",
		sess.LocalPort(), conn.RemoteAddr(), conn.LocalAddr())
}

func (ms *MsocksServer) Serve(listener net.Listener) (err error) {
	var conn net.Conn

	for {
		conn, err = listener.Accept()
		if err != nil {
			log.Error("%s", err)
			continue
		}
		go func() {
			defer conn.Close()
			ms.Handler(conn)
		}()
	}
	return
}
