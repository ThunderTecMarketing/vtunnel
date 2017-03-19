package msocks

import (
	"github.com/FTwOoO/vpncore/net/conn/message/msgpack"
	"fmt"
	"reflect"
)

func init() {
	msgpack.RegisterMessage(FrameSynResult{}.Cmd(), reflect.TypeOf(FrameSynResult{}))
	msgpack.RegisterMessage(FrameData{}.Cmd(), reflect.TypeOf(FrameData{}))
	msgpack.RegisterMessage(FrameSyn{}.Cmd(), reflect.TypeOf(FrameSyn{}))
	msgpack.RegisterMessage(FrameFin{}.Cmd(), reflect.TypeOf(FrameFin{}))
	msgpack.RegisterMessage(FrameRst{}.Cmd(), reflect.TypeOf(FrameRst{}))
	msgpack.RegisterMessage(FrameDns{}.Cmd(), reflect.TypeOf(FrameDns{}))
}

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

type ConnInfo struct {
	Network string
	SrcHost string
	SrcPort uint16
	DstHost string
	DstPort uint16
}

func (c *ConnInfo) String() (s string) {
	return fmt.Sprintf("%s [%s:%d]->[%s:%d]", c.Network, c.SrcHost, c.SrcPort, c.DstHost, c.DstPort)
}

type Frame interface {
	msgpack.Message
	GetStreamId() uint16
}

type FrameReceiver interface {
	ReceiveFrame(Frame) error
	CloseFrame() error
}

type FrameSender interface {
	SendFrame(Frame) error
	CloseFrame() error
}


//go:generate msgp
//msgp:tuple ConnInfo FrameBase FrameResult FrameData FrameSyn FrameWnd FrameFin FrameRst FramePing FrameDns FrameSpam


type FrameSynResult struct {
	StreamId uint16
	Errno    uint32
}

func (f FrameSynResult) GetStreamId() uint16 {
	return f.StreamId
}

func (z FrameSynResult) Cmd() msgpack.MessageType {
	return MSG_RESULT
}

type FrameData struct {
	StreamId uint16
	Data     []byte
}

func (f FrameData) GetStreamId() uint16 {
	return f.StreamId
}

func (z FrameData) Cmd() msgpack.MessageType {
	return MSG_DATA
}

type FrameSyn struct {
	StreamId uint16
	Address  ConnInfo
}

func (f FrameSyn) GetStreamId() uint16 {
	return f.StreamId
}

func (z FrameSyn) Cmd() msgpack.MessageType {
	return MSG_SYN
}

type FrameFin struct {
	StreamId uint16
}

func (f FrameFin) GetStreamId() uint16 {
	return f.StreamId
}

func (z FrameFin) Cmd() msgpack.MessageType {
	return MSG_FIN
}

type FrameRst struct {
	StreamId uint16
}

func (f FrameRst) GetStreamId() uint16 {
	return f.StreamId
}

func (z FrameRst) Cmd() msgpack.MessageType {
	return MSG_RST
}

type FrameDns struct {
	StreamId uint16
	Data     []byte
}

func (f FrameDns) GetStreamId() uint16 {
	return f.StreamId
}

func (z FrameDns) Cmd() msgpack.MessageType {
	return MSG_DNS
}
