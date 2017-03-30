package msocks

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"
	"github.com/FTwOoO/vtunnel/util"
)

const (
	ST_UNKNOWN = iota
	ST_SYN_RECV
	ST_SYN_SENT
	ST_EST
	ST_FIN
)


type Stream interface {
	net.Conn
	FrameReceiver
	String() string
}

var _ Stream = new(Conn)


//implements Stream interface
type Conn struct {
	Address     ConnInfo

	statusLock  sync.Mutex
	status      uint8

	session     *Session
	streamId    uint16

	chSynResult chan SyncResultCode

	rlock       sync.Mutex
	wlock       sync.Mutex

	rqueue      *util.Queue
}

func NewConn(status uint8, streamid uint16, session *Session, address ConnInfo) (c *Conn) {
	c = &Conn{
		status:   status,
		session:  session,
		streamId: streamid,
		Address:  address,
		rqueue:   util.NewQueue(),
	}
	return
}

func (c *Conn) GetStreamId() uint16 {
	return c.streamId
}

func (c *Conn) String() (s string) {
	return fmt.Sprintf("Stream[ID:%d %s]", c.streamId, c.Address.String())
}

func (c *Conn) Read(data []byte) (n int, err error) {
	var v interface{}
	c.rlock.Lock()
	defer c.rlock.Unlock()

	target := data[:]
	block := true

	var r_rest = []byte(nil)

	for len(target) > 0 {
		if r_rest == nil {
			// reader should be blocked in here
			v, err = c.rqueue.Pop(block)
			if err == util.ErrQueueClosed {
				err = io.EOF
			}
			if err != nil {
				return
			}
			if v == nil {
				break
			}
			r_rest = v.([]byte)
		}

		size := copy(target, r_rest)
		target = target[size:]
		n += size
		block = false

		if len(r_rest) > size {
			r_rest = r_rest[size:]
		} else {
			// take all data in rest
			r_rest = nil
		}
	}

	if r_rest != nil {
		c.rqueue.PushFront(r_rest)
	}

	return
}

func (c *Conn) Write(data []byte) (n int, err error) {
	c.wlock.Lock()
	defer c.wlock.Unlock()

	for len(data) > 0 {
		size := uint32(len(data))

		//limit size < 4Kb
		switch {
		case size > 8 * 1024:
			size = uint32(3 * 1024 + rand.Intn(1024))
		case 4 * 1024 < size && size <= 8 * 1024:
			size /= 2
		}

		err = c.writeSlice(data[:size])

		if err != nil {
			log.Errorf("%s", err)
			return
		}
		log.Debugf("%s send chunk %d bytes", c.String(), size)

		data = data[size:]
		n += int(size)
	}
	log.Infof("%s sent %d bytes.", c.String(), n)
	return
}

func (c *Conn) writeSlice(data []byte) (err error) {
	f := &FrameData{StreamId:c.streamId, Data:data}

	if c.status != ST_EST {
		log.Errorf("status %d found in write slice", c.status)
		return ErrState
	}

	err = c.session.SendFrame(f)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}


func (c *Conn) sendFin() (err error) {
	fb := &FrameFin{StreamId: c.streamId}
	err = c.session.SendFrame(fb)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	return
}


func (c *Conn) LocalAddr() net.Addr {
	return &Addr{
		NetworkType:c.Address.Network,
		Address:fmt.Sprintf("%s:%d", c.Address.SrcHost, c.Address.SrcPort),
		Streamid:c.streamId,
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	return &Addr{
		NetworkType:c.Address.Network,
		Address:fmt.Sprintf("%s:%d", c.Address.DstHost, c.Address.DstPort),
		Streamid:c.streamId,
	}
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *Conn) GetStatus() (st string) {
	switch c.status {
	case ST_SYN_RECV:
		return "SYN_RECV"
	case ST_SYN_SENT:
		return "SYN_SENT"
	case ST_EST:
		return "ESTAB"
	case ST_FIN:
		return "FIN_WAIT"
	}
	return "UNKNOWN"
}

type Addr struct {
	NetworkType string
	Address     string
	Streamid    uint16
}

func (a *Addr) String() string {
	return a.Address
}

func (a *Addr) Network() string {
	return a.NetworkType
}
