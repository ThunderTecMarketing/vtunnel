package msocks

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"github.com/miekg/dns"
	"github.com/FTwOoO/vpncore/net/conn"

	mdns "github.com/FTwOoO/vpncore/net/dns"
	"github.com/FTwOoO/vtunnel/util"
)

type Session struct {
	conn        conn.ObjectIO
	next_id     uint16

	streamsLock sync.Mutex
	streams     map[uint16]Stream

	dialer      Dialer
	dnsServer   *mdns.DNSServer
}

func NewSession(conn conn.ObjectIO, dnsServer *mdns.DNSServer) (s *Session) {
	s = &Session{
		conn:     conn,
		dnsServer: dnsServer,
		streams:    make(map[uint16]Stream, 0),
	}
	log.Noticef("session %s created.", s.String())
	return
}

func (s *Session) String() string {
	return "Session()"
}

func (s *Session) Close() (err error) {
	log.Warningf("%s close all streams,  %d streams closed", s.String(), len(s.streams))
	defer s.conn.Close()
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	for _, v := range s.streams {
		v.CloseFrame()
	}
	return
}

func (s *Session) Run() {
	defer s.Close()

	for {
		obj, err := s.conn.Read()
		if err != nil {
			log.Error(err)
			return
		}

		var f Frame
		var ok bool

		if f, ok = obj.(Frame); !ok {
			log.Error("Receive object that is not Frame")
			return
		}

		switch ft := f.(type) {
		default:
			log.Errorf("%s", ErrUnexpectedPkg.Error())
			return
		case *FrameSynResult, *FrameData, *FrameFin, *FrameRst:
			err = s.on_stream_packet(f)
			if err != nil {
				log.Errorf("%s send to stream[%d] failed, err: %s.", s.String(), f.GetStreamId(), err.Error())
				continue
			}
		case *FrameSyn:
			err = s.on_syn(ft)
			if err != nil {
				log.Errorf("syn failed: %s", err.Error())
				return
			}
		case *FrameDns:
			err = s.on_dns(ft)
			if err != nil {
				log.Errorf("dns failed: %s", err.Error())
				return
			}

		}
	}
}

func (s *Session) on_stream_packet(f Frame) (err error) {
	streamid := f.GetStreamId()
	c, err := s.GetStreamById(streamid)
	if err != nil {
		return err
	}

	err = c.ReceiveFrame(f)
	if err != nil {
		return
	}
	return nil
}

func (s *Session) on_syn(ft *FrameSyn) (err error) {
	// lock streamid temporary, with status sync recved
	c := NewConn(ST_SYN_RECV, ft.GetStreamId(), s, ft.Address)

	err = s.PutStreamIntoId(ft.GetStreamId(), c)
	if err != nil {
		log.Errorf("%s", err)

		fb := &FrameSynResult{StreamId:ft.GetStreamId(), Errno:ERR_IDEXIST}
		err := s.SendFrame(fb)
		if err != nil {
			return err
		}
		return nil
	}

	// it may toke long time to connect with target address
	// so we use goroutine to return back loop
	go func() {
		var err error
		var connection net.Conn

		var network = ft.Address.Network
		var address = fmt.Sprintf("%s:%d", ft.Address.DstHost, ft.Address.DstPort)
		log.Debugf("try to connect %s => %s:%s.", c.String(), network, address)

		if dialer, ok := s.dialer.(*TcpDialer); ok {
			connection, err = dialer.DialTimeout(network, address, DIAL_TIMEOUT * time.Second)
		} else {
			connection, err = s.dialer.Dial(network, address)
		}

		if err != nil {
			log.Errorf("%s", err)
			fb := &FrameSynResult{StreamId:ft.GetStreamId(), Errno:ERR_CONNFAILED}
			err = s.SendFrame(fb)
			if err != nil {
				log.Errorf("%s", err)
			}
			c.CloseFrame()
			return
		}

		fb := &FrameSynResult{StreamId:ft.GetStreamId(), Errno:ERR_NONE}
		err = s.SendFrame(fb)
		if err != nil {
			log.Errorf("%s", err)
			return
		}
		c.status = ST_EST

		go util.CopyLink(connection, c)
		return
	}()
	return
}

func (s *Session) writeDNS(ctx mdns.QueryContext, m []byte) (err error) {
	streamId, ok := ctx.(uint16)
	if !ok {
		return errors.New("Unexpected context")
	}

	fr := &FrameDns{StreamId:streamId, Data:m}
	err = s.SendFrame(fr)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) on_dns(ft *FrameDns) (err error) {
	req := new(dns.Msg)
	err = req.Unpack(ft.Data)
	if err != nil {
		return ErrDnsMsgIllegal
	}

	s.dnsServer.QueryByDNSMsg(req, ft.GetStreamId(), mdns.RawHandler(s.writeDNS))
	return
}

