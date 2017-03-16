package msocks

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/FTwOoO/vpncore/net/conn"
	"strconv"
)

type Session struct {
	conn        conn.ObjectIO
	next_id     uint16

	streamsLock sync.Mutex
	streams     map[uint16]FrameSender

	dialer      Dialer
}

func NewSession(conn conn.ObjectIO) (s *Session) {
	s = &Session{
		conn:     conn,
		streams:    make(map[uint16]FrameSender, 0),
	}
	log.Notice("session %s created.", s.String())
	return
}

func (s *Session) String() string {
	return fmt.Sprintf("%d", s.LocalPort())
}

func (s *Session) GetSize() int {
	return len(s.streams)
}

func (s *Session) GetStreamById(id uint16) (FrameSender, error) {
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	c, ok := s.streams[id]
	if !ok || c == nil {
		return nil, ErrStreamNotExist
	}
	return c, nil
}

func (s *Session) GetStreams() (ports []*Conn) {
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	for _, fs := range s.streams {
		if c, ok := fs.(*Conn); ok {
			ports = append(ports, c)
		}
	}
	return
}

func (s *Session) RemoveStream(streamid uint16) (err error) {
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	_, ok := s.streams[streamid]
	if !ok {
		return fmt.Errorf("streamid(%d) not exist.", streamid)
	}
	delete(s.streams, streamid)
	log.Info("%s remove port %d.", s.String(), streamid)
	return
}

type ConnSlice []*Conn

func (cs ConnSlice) Len() int           { return len(cs) }
func (cs ConnSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs ConnSlice) Less(i, j int) bool { return cs[i].streamId < cs[j].streamId
}

func (s *Session) GetSortedStreams() (ports ConnSlice) {
	ports = s.GetStreams()
	sort.Sort(ports)
	return
}

func (s *Session) PutIntoNextId(fs FrameSender) (id uint16, err error) {
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	startid := s.next_id
	for  {

		_, ok := s.streams[s.next_id]
		if !ok {
			break
		}

		s.next_id += 1
		if s.next_id == startid {
			err = errors.New("run out of stream id")
			log.Error("%s", err)
			return
		}
	}
	id = s.next_id
	s.next_id += 1
	log.Debug("%s put into next id %d: %p.", s.String(), id, fs)

	s.streams[id] = fs
	return
}

func (s *Session) PutIntoId(id uint16, fs FrameSender) (err error) {
	log.Debug("%s put into id %d: %p.", s.String(), id, fs)
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	_, ok := s.streams[id]
	if ok {
		return ErrIdExist
	}

	s.streams[id] = fs
	return
}



func (s *Session) Close() (err error) {
	log.Warning("close all connects (%d) for session: %s.",
		len(s.streams), s.String())
	defer s.conn.Close()
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	for _, v := range s.streams {
		v.CloseFrame()
	}
	return
}

func (s *Session) LocalAddr() net.Addr {
	return nil
}

func (s *Session) RemoteAddr() net.Addr {
	return nil
}

func (s *Session) LocalPort() int {
	addr, ok := s.LocalAddr().(*net.TCPAddr)
	if !ok {
		return -1
	}
	return addr.Port
}

func (s *Session) SendFrame(f Frame) (err error) {
	err = s.conn.Write(f)
	if err != nil {
		return
	}

	return
}

func (s *Session) Run() {
	defer s.Close()

	for {
		f, err := ReadFrame(s.conn)
		if err != nil {
			log.Error("%s", err)
			return
		}

		switch ft := f.(type) {
		default:
			log.Error("%s", ErrUnexpectedPkg.Error())
			return
		case *FrameSynResult, *FrameData, *FrameFin, *FrameRst:
			err = s.sendFrameToStream(f)
			if err != nil {
				log.Error("%s(%d) send failed, err: %s.",
					s.String(), f.GetStreamid(), err.Error())
				return
			}
		case *FrameSyn:
			err = s.on_syn(ft)
			if err != nil {
				log.Error("syn failed: %s", err.Error())
				return
			}
		case *FrameDns:
			err = s.on_dns(ft)
			if err != nil {
				log.Error("dns failed: %s", err.Error())
				return
			}
		case *FramePing:
		case *FrameSpam:
		}
	}
}

func (s *Session) sendFrameToStream(f Frame) (err error) {
	streamid := f.GetStreamid()
	c, err := s.GetStreamById(streamid)
	if err != nil {
		return err
	}

	err = c.SendFrame(f)
	if err != nil {
		return
	}
	return nil
}


func (s *Session) on_syn(ft *FrameSyn) (err error) {
	// lock streamid temporary, with status sync recved
	c := NewConn(ST_SYN_RECV, ft.Streamid, s, ft.Address)

	err = s.PutIntoId(ft.Streamid, c)
	if err != nil {
		log.Error("%s", err)

		fb := &FrameSynResult{FrameBase.Streamid:ft.Streamid, Errno:ERR_IDEXIST}
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
		var conn net.Conn
		log.Debug("try to connect %s => %s:%s.", c.String(), ft.Network, ft.Address)

		if dialer, ok := s.dialer.(*TcpDialer); ok {
			conn, err = dialer.DialTimeout(ft.Network, ft.Address, DIAL_TIMEOUT*time.Second)
		} else {
			conn, err = s.dialer.Dial(ft.Network, ft.Address)
		}

		if err != nil {
			log.Error("%s", err)
			fb := &FrameSynResult{FrameBase.Streamid:ft.Streamid, Errno:ERR_CONNFAILED}
			err = s.SendFrame(fb)
			if err != nil {
				log.Error("%s", err)
			}
			c.Final()
			return
		}

		fb := &FrameSynResult{FrameBase.Streamid:ft.Streamid, Errno:ERR_NONE}
		err = s.SendFrame(fb)
		if err != nil {
			log.Error("%s", err)
			return
		}
		c.status = ST_EST

		go CopyLink(conn, c)
		log.Notice("connected %s => %s:%s.", c.String(), ft.Network, ft.Address)
		return
	}()
	return
}


func (s *Session) on_dns(ft *FrameDns) (err error) {
	req := new(dns.Msg)
	err = req.Unpack(ft.Data)
	if err != nil {
		return ErrDnsMsgIllegal
	}

	if req.Response {
		// ignore send fail, maybe just timeout.
		// should I log this ?
		return s.sendFrameToStream(ft)
	}

	log.Info("dns query for %s.", req.Question[0].Name)

	d, ok := DefaultLookuper.(*DnsLookup)
	if !ok {
		return ErrNoDnsServer
	}
	res, err := d.Exchange(req)
	if err != nil {
		log.Error("%s", err.Error())
		return nil
	}

	if DEBUGDNS {
		DebugDNS(res, req.Question[0].Name)
	}

	// send response back from streamid
	b, err := res.Pack()
	if err != nil {
		log.Error("%s", ErrDnsMsgIllegal.Error())
		return nil
	}

	fr := &FrameDns{FrameBase.Streamid:ft.GetStreamid(), Data:b}
	err = s.SendFrame(fr)
	return
}


func (s *Session) Dial(network, address string) (c *Conn, err error) {
	dst, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return
	}

	c = NewConn(ST_SYN_SENT, 0, s, ConnInfo{Network:network, DstHost:dst, DstPort:port})
	streamid, err := s.PutIntoNextId(c)
	if err != nil {
		return
	}
	c.streamId = streamid

	//log.Info("try dial %s => %s.", s.conn.RemoteAddr().String(), address)
	err = c.WaitForConn()
	if err != nil {
		return
	}

	return c, nil
}

/*
func (s *Session) LookupIP(host string) (addrs []net.IP, err error) {
	ip := net.ParseIP(host)
	if ip != nil {
		return []net.IP{ip}, nil
	}

	cfs := CreateChanFrameSender(0)
	streamid, err := s.PutIntoNextId(&cfs)
	if err != nil {
		return
	}
	defer func() {
		err := s.RemoveStream(streamid)
		if err != nil {
			log.Error("%s", err.Error())
		}
	}()

	req, freq, err := MakeDnsFrame(host, dns.TypeA, streamid)
	if err != nil {
		return
	}

	err = s.SendFrame(freq)
	if err != nil {
		return
	}

	fres, err := cfs.RecvWithTimeout(DNS_TIMEOUT * time.Second)
	if err != nil {
		return
	}

	addrs, err = ParseDnsFrame(fres, req)
	return
}
*/
