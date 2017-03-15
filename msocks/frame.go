package msocks

import (
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"github.com/FTwOoO/vpncore/net/conn"
	"errors"
)

const (
	MSG_UNKNOWN msgpack.MessageType = 1
	MSG_RESULT msgpack.MessageType = 2
	MSG_AUTH msgpack.MessageType = 3
	MSG_DATA msgpack.MessageType = 4
	MSG_SYN msgpack.MessageType = 5
	MSG_FIN msgpack.MessageType = 7
	MSG_RST msgpack.MessageType = 8
	MSG_PING msgpack.MessageType = 9
	MSG_DNS msgpack.MessageType = 10
	MSG_SPAM msgpack.MessageType = 11
)

type Frame interface {
	GetStreamid() uint16
	msgpack.Message
}

type FrameSender interface {
	SendFrame(Frame) error
	CloseFrame() error
}

func ReadFrame(r conn.ObjectIO) (Frame, error) {
	obj, err :=  r.Read()
	if err != nil {
		return nil, err
	}

	var f Frame
	var ok bool

	if f, ok = obj.(Frame); !ok {
		return nil, errors.New("Receive object that is not Frame")
	}

	return f, nil
}

//go:generate msgp
//msgp:tuple FrameBase FrameResult FrameData FrameSyn FrameWnd FrameFin FrameRst FramePing FrameDns FrameSpam


type FrameBase struct {
	Streamid uint16
}

func (f FrameBase) GetStreamid() uint16 {
	return f.Streamid
}

type FrameSynResult struct {
	*FrameBase
	Errno uint32
}

func (z FrameSynResult) Cmd() msgpack.MessageType {
	return MSG_RESULT
}

type FrameData struct {
	*FrameBase
	Data []byte
}

func (z FrameData) Cmd() msgpack.MessageType {
	return MSG_DATA
}

type FrameSyn struct {
	*FrameBase
	Network string
	Address string
}

func (z FrameSyn) Cmd() msgpack.MessageType {
	return MSG_SYN
}

type FrameFin struct {
	*FrameBase
}

func (z FrameFin) Cmd() msgpack.MessageType {
	return MSG_FIN
}

type FrameRst struct {
	*FrameBase
}

func (z FrameRst) Cmd() msgpack.MessageType {
	return MSG_RST
}

type FramePing struct {
	*FrameBase
}

func (z FramePing) Cmd() msgpack.MessageType {
	return MSG_PING
}

type FrameDns struct {
	*FrameBase
	Data []byte
}

func (z FrameDns) Cmd() msgpack.MessageType {
	return MSG_DNS
}

type FrameSpam struct {
	*FrameBase
	Data []byte
}

func (z FrameSpam) Cmd() msgpack.MessageType {
	return MSG_SPAM
}