package msocks

import (
	"errors"
	"math/rand"
	"time"
	"github.com/FTwOoO/go-logger"
)

const (
	DIAL_RETRY = 2
	DIAL_TIMEOUT = 5
	DNS_TIMEOUT = 30
)

var (
	ErrNoSession = errors.New("session in pool but can't pick one.")
	ErrSessionNotFound = errors.New("session not found.")
	ErrAuthFailed = errors.New("auth failed.")
	ErrAuthTimeout = errors.New("auth timeout %s.")
	ErrStreamNotExist = errors.New("stream not exist.")
	ErrUnexpectedPkg = errors.New("unexpected package.")
	ErrNotSyn = errors.New("frame result in conn which status is not syn.")
	ErrIdExist = errors.New("frame sync stream id exist.")
	ErrState = errors.New("status error.")
	ErrUnknownState = errors.New("unknown status.")
	ErrChanClosed = errors.New("chan closed.")
	ErrDnsTimeOut = errors.New("dns timeout.")
	ErrDnsMsgIllegal = errors.New("dns message illegal.")
)

var log *logger.Logger

func RegisterLogger(log_file string, log_level string) (err error) {
	log, err = logger.NewLogger(log_file, log_level)
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
