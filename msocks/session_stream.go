package msocks

import (
	"sort"
	"fmt"
	"errors"
)

func (s *Session) GetStreamsSize() int {
	return len(s.streams)
}

func (s *Session) GetStreamById(id uint16) (Stream, error) {
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

func (s *Session) GetSortedStreams() (ports ConnSlice) {
	ports = s.GetStreams()
	sort.Sort(ports)
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
	log.Infof("%s remove stream[%d].", s.String(), streamid)
	return
}

type ConnSlice []*Conn

func (cs ConnSlice) Len() int {
	return len(cs)
}
func (cs ConnSlice) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
func (cs ConnSlice) Less(i, j int) bool {
	return cs[i].streamId < cs[j].streamId
}

func (s *Session) PutStreamIntoNextId(fs Stream) (id uint16, err error) {
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	startid := s.next_id
	for {
		_, ok := s.streams[s.next_id]
		if !ok {
			break
		}

		s.next_id += 1
		if s.next_id == startid {
			err = errors.New("run out of stream id")
			log.Errorf("%s", err)
			return
		}
	}
	id = s.next_id
	s.next_id += 1
	log.Debugf("%s put into next id %d: %s.", s.String(), id, fs)

	s.streams[id] = fs
	return
}

func (s *Session) PutStreamIntoId(id uint16, fs Stream) (err error) {
	log.Debugf("%s put %s into id %d", s.String(), fs.String(), id)
	s.streamsLock.Lock()
	defer s.streamsLock.Unlock()

	_, ok := s.streams[id]
	if ok {
		return ErrIdExist
	}

	s.streams[id] = fs
	return
}