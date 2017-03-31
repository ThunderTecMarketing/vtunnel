package msocks

import (
	"math/rand"
	"sync"
)

type ClientSessionPool struct {
	sessionsMu     sync.Mutex
	sessions       map[*Session]struct{}

	factoryMu      sync.Mutex
	SessionDialers []ObjectDialer

	MinSess        int
	MaxConn        int
}

func NewClientSessionPool(MinSess, MaxConn int, dialers []ObjectDialer) (sp *ClientSessionPool) {
	if MinSess == 0 {
		MinSess = 1
	}
	if MaxConn == 0 {
		MaxConn = 16
	}
	sp = &ClientSessionPool{
		sessions:    make(map[*Session]struct{}, 0),
		MinSess: MinSess,
		MaxConn: MaxConn,
		SessionDialers:dialers,
	}
	return
}

func (sp *ClientSessionPool) AddSessionDialer(sd ObjectDialer) {
	sp.factoryMu.Lock()
	defer sp.factoryMu.Unlock()
	sp.SessionDialers = append(sp.SessionDialers, sd)
}

func (sp *ClientSessionPool) CleanSessions() {
	sp.sessionsMu.Lock()
	defer sp.sessionsMu.Unlock()
	for s, _ := range sp.sessions {
		s.Close()
	}
	sp.sessions = make(map[*Session]struct{}, 0)
}

func (sp *ClientSessionPool) GetSessionsCount() int {
	return len(sp.sessions)
}

func (sp *ClientSessionPool) GetSessions() (sess map[*Session]struct{}) {
	return sp.sessions
}

func (sp *ClientSessionPool) Add(s *Session) {
	sp.sessionsMu.Lock()
	defer sp.sessionsMu.Unlock()
	sp.sessions[s] = struct{}{}
}

func (sp *ClientSessionPool) Remove(s *Session) (err error) {
	sp.sessionsMu.Lock()
	defer sp.sessionsMu.Unlock()
	if _, ok := sp.sessions[s]; !ok {
		return ErrSessionNotFound
	}
	delete(sp.sessions, s)
	return
}

func (sp *ClientSessionPool) Get() (sess *Session, err error) {
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

func (sp *ClientSessionPool) createSession(checker func() bool) (err error) {
	sp.factoryMu.Lock()
	defer sp.factoryMu.Unlock()

	if checker != nil && !checker() {
		return
	}

	var sess *Session

	start := rand.Int()
	end := start + DIAL_RETRY * len(sp.SessionDialers)
	for i := start; i < end; i++ {
		asf := sp.SessionDialers[i % len(sp.SessionDialers)]

		conn, err := asf.Dial()
		if err != nil {
			log.Errorf("%s", err)
			continue
		}

		sess = NewSession(conn, nil, nil)
		break
	}

	if err != nil {
		log.Errorf("can't connect to any server, quit.")
		return
	}

	sp.Add(sess)
	go sp.RunSession(sess)
	return
}

func (sp *ClientSessionPool) getLessUsedSession() (sess *Session, size int) {
	size = -1
	for s, _ := range sp.sessions {
		if size == -1 || s.GetStreamsSize() < size {
			sess = s
			size = s.GetStreamsSize()
		}
	}
	return
}

func (sp *ClientSessionPool) RunSession(sess *Session) {
	defer func() {
		err := sp.Remove(sess)
		if err != nil {
			log.Errorf("%s", err)
			return
		}
	}()

	sess.Run()
	return
}
