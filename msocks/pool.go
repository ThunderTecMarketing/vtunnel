package msocks

import (
	"math/rand"
	"net"
	"sync"
)

type ClientDialerFactory struct {
	ObjectDialer

}

func (sf *ClientDialerFactory) CreateSession() (s *Session, err error) {
	conn, err := sf.ObjectDialer.Dial()
	if err != nil {
		return
	}

	s = NewSession(conn)
	return
}

type SessionPool struct {
	sessionsMu       sync.Mutex
	factoryMu        sync.Mutex
	sessions         map[*Session]struct{}
	sessionFactories []*ClientDialerFactory
	MinSess          int
	MaxConn          int
}

func CreateSessionPool(MinSess, MaxConn int) (sp *SessionPool) {
	if MinSess == 0 {
		MinSess = 1
	}
	if MaxConn == 0 {
		MaxConn = 16
	}
	sp = &SessionPool{
		sessions:    make(map[*Session]struct{}, 0),
		MinSess: MinSess,
		MaxConn: MaxConn,
	}
	return
}

func (sp *SessionPool) AddSessionFactory(sf ClientDialerFactory) {
	sp.factoryMu.Lock()
	defer sp.factoryMu.Unlock()
	sp.sessionFactories = append(sp.sessionFactories, sf)
}

func (sp *SessionPool) CleanSessions() {
	sp.sessionsMu.Lock()
	defer sp.sessionsMu.Unlock()
	for s, _ := range sp.sessions {
		s.Close()
	}
	sp.sessions = make(map[*Session]struct{}, 0)
}

func (sp *SessionPool) GetSize() int {
	return len(sp.sessions)
}

func (sp *SessionPool) GetSessions() (sess map[*Session]struct{}) {
	return sp.sessions
}

func (sp *SessionPool) Add(s *Session) {
	sp.sessionsMu.Lock()
	defer sp.sessionsMu.Unlock()
	sp.sessions[s] = struct{}{}
}

func (sp *SessionPool) Remove(s *Session) (err error) {
	sp.sessionsMu.Lock()
	defer sp.sessionsMu.Unlock()
	if _, ok := sp.sessions[s]; !ok {
		return ErrSessionNotFound
	}
	delete(sp.sessions, s)
	return
}

func (sp *SessionPool) Get() (sess *Session, err error) {
	if len(sp.sessions) == 0 {
		err = sp.createSession(func() bool {
			return len(sp.sessions) == 0
		})
		if err != nil {
			return nil, err
		}
	}

	sess, size := sp.getLessUsedSession()
	if sess == nil {
		return nil, ErrNoSession
	}

	if size > sp.MaxConn || len(sp.sessions) < sp.MinSess {
		go sp.createSession(func() bool {
			if len(sp.sessions) < sp.MinSess {
				return true
			}
			// normally, size == -1 should never happen
			_, size := sp.getLessUsedSession()
			return size > sp.MaxConn
		})
	}
	return
}

// Randomly select a server, try to connect with it. If it is failed, try next.
// Repeat for DIAL_RETRY times.
// Each time it will take 2 ^ (net.ipv4.tcp_syn_retries + 1) - 1 second(s).
// eg. net.ipv4.tcp_syn_retries = 4, connect will timeout in 2 ^ (4 + 1) -1 = 31s.
func (sp *SessionPool) createSession(checker func() bool) (err error) {
	sp.factoryMu.Lock()
	defer sp.factoryMu.Unlock()

	if checker != nil && !checker() {
		return
	}

	var sess *Session

	start := rand.Int()
	end := start + DIAL_RETRY*len(sp.sessionFactories)
	for i := start; i < end; i++ {
		asf := sp.sessionFactories[i%len(sp.sessionFactories)]
		sess, err = asf.CreateSession()
		if err != nil {
			log.Error("%s", err)
			continue
		}
		break
	}

	if err != nil {
		log.Critical("can't connect to any server, quit.")
		return
	}
	log.Notice("session created.")

	sp.Add(sess)
	go sp.sessRun(sess)
	return
}

func (sp *SessionPool) getLessUsedSession() (sess *Session, size int) {
	size = -1
	for s, _ := range sp.sessions {
		if size == -1 || s.GetSize() < size {
			sess = s
			size = s.GetSize()
		}
	}
	return
}

func (sp *SessionPool) sessRun(sess *Session) {
	defer func() {
		err := sp.Remove(sess)
		if err != nil {
			log.Error("%s", err)
			return
		}
	}()

	sess.Run()
	log.Warning("session runtime quit, reboot from connect.")
	return
}

func (sp *SessionPool) Dial(network, address string) (net.Conn, error) {
	sess, err := sp.Get()
	if err != nil {
		return nil, nil
	}
	return sess.Dial(network, address)
}

/*
func (sp *SessionPool) LookupIP(host string) (addrs []net.IP, err error) {
	sess, err := sp.Get()
	if err != nil {
		return
	}
	return sess.LookupIP(host)
}
*/
