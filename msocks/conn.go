package msocks

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	ST_UNKNOWN = iota
	ST_SYN_RECV
	ST_SYN_SENT
	ST_EST
	ST_FIN
)

type ConnInfo struct {
	Network string
	SrcIp string
	SrcPort uint16
	DstHost string
	DstPort uint16
}

type Conn struct {
	//The target
	Address     ConnInfo

	statusLock  sync.Mutex
	status      uint8

	session     *Session
	streamId    uint16
	sender      FrameSender

	chSynResult chan uint32

	rlock       sync.Mutex // this should used to block reader and reader, not writer
	wlock       sync.Mutex

	rqueue      *Queue
}

func NewConn(status uint8, streamid uint16, session *Session, address ConnInfo) (c *Conn) {
	c = &Conn{
		status:   status,
		session:  session,
		streamId: streamid,
		sender:   session,
		Address:  address,
		rqueue:   NewQueue(),
	}
	return
}

func (c *Conn) GetStreamId() uint16 {
	return c.streamId
}

func (c *Conn) GetAddress() (s string) {
	return fmt.Sprintf("%s:%s", c.Network, c.Address)
}

func (c *Conn) String() (s string) {
	return fmt.Sprintf("%d(%d)", c.session.LocalPort(), c.streamId)
}

func recvWithTimeout(ch chan uint32, t time.Duration) (errno uint32) {
	var ok bool
	ch_timeout := time.After(t)
	select {
	case errno, ok = <-ch:
		if !ok {
			return ERR_CLOSED
		}
	case <-ch_timeout:
		return ERR_TIMEOUT
	}
	return
}

func (c *Conn) WaitForConn() (err error) {
	c.chSynResult = make(chan uint32, 0)

	fb := &FrameSyn{c.streamId, c.Network, c.Address}
	err = c.session.SendFrame(fb)
	if err != nil {
		log.Error("%s", err)
		c.Final()
		return
	}

	errno := recvWithTimeout(c.chSynResult, DIAL_TIMEOUT * time.Second)
	if errno != ERR_NONE {
		log.Error("remote connect %s failed for %d.", c.String(), errno)
		c.Final()
	} else {
		log.Notice("%s connected: %s => %s.", c.Network, c.String(), c.Address)
	}

	c.chSynResult = nil
	return
}

func (c *Conn) Final() {
	c.rqueue.Close()

	err := c.session.RemoveStream(c.streamId)
	if err != nil {
		log.Error("%s", err)
	}

	log.Notice("%s final.", c.String())
	c.status = ST_UNKNOWN
	return
}

func (c *Conn) Close() (err error) {
	log.Info("close %s.", c.String())
	c.statusLock.Lock()
	defer c.statusLock.Unlock()

	fb := &FrameFin{FrameBase.Streamid: c.streamId}
	err = c.sender.SendFrame(fb)
	if err != nil {
		log.Error("%s", err)
		return
	}
	c.Final()
	return
}

func (c *Conn) SendFrame(f Frame) (err error) {
	switch ft := f.(type) {
	default:
		err = ErrUnexpectedPkg
		log.Error("%s", err)
		c.Close()
		return
	case *FrameSynResult:
		return c.InSynResult(ft.Errno)
	case *FrameData:
		return c.InData(ft)
	case *FrameFin:
		return c.InFin(ft)
	case *FrameRst:
		log.Debug("reset %s.", c.String())
		c.Final()
	}
	return
}

func (c *Conn) InSynResult(errno uint32) (err error) {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()

	if c.status != ST_SYN_SENT {
		return ErrNotSyn
	}

	if errno == ERR_NONE {
		c.status = ST_EST
	} else {
		c.Final()
	}

	select {
	case c.chSynResult <- errno:
	default:
	}
	return
}

func (c *Conn) InData(ft *FrameData) (err error) {
	log.Info("%s recved %d bytes.", c.String(), len(ft.Data))
	err = c.rqueue.Push(ft.Data)
	if err != nil {
		return
	}
	return
}

func (c *Conn) InFin(ft *FrameFin) (err error) {
	// always need to close read pipe
	// coz fin means remote will never send data anymore
	c.rqueue.Close()

	c.statusLock.Lock()
	defer c.statusLock.Unlock()

	c.Final()
	return

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
			if err == ErrQueueClosed {
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

		err = c.WriteSlice(data[:size])

		if err != nil {
			log.Error("%s", err)
			return
		}
		log.Debug("%s send chunk size %d at %d.", c.String(), size, n)

		data = data[size:]
		n += int(size)
	}
	log.Info("%s sent %d bytes.", c.String(), n)
	return
}

func (c *Conn) WriteSlice(data []byte) (err error) {
	f := &FrameData{FrameBase.Streamid:c.streamId, Data:data}

	if c.status != ST_EST {
		log.Error("status %d found in write slice", c.status)
		return ErrState
	}

	err = c.sender.SendFrame(f)
	if err != nil {
		log.Error("%s", err)
		return
	}
	return
}

func (c *Conn) LocalAddr() net.Addr {
	return &Addr{
		Network:c.Address.Network,
		Address:fmt.Sprintf("%s:%d", c.Address.SrcIp, c.Address.SrcPort),
		c.streamId,
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	return &Addr{
		Network:c.Address.Network,
		Address:fmt.Sprintf("%s:%d", c.Address.DstHost, c.Address.DstPort),
		c.streamId,
	}
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

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

type Addr struct {
	Network string
	Address string
	streamid uint16
}

func (a *Addr) String() string {
	return a.Address
}

func (a *Addr) Network() string {
	return a.Network
}
