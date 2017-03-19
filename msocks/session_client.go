package msocks

import (
	"strconv"
	"net"
)

func (s *Session) Dial(network, address string) (c *Conn, err error) {
	dst, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return
	}

	c = NewConn(ST_SYN_SENT, 0, s, ConnInfo{Network:network, DstHost:dst, DstPort:uint16(port)})
	streamid, err := s.PutStreamIntoNextId(c)
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
