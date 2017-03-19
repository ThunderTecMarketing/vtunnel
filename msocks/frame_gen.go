package msocks

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ConnInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var zxvk uint32
	zxvk, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zxvk != 5 {
		err = msgp.ArrayError{Wanted: 5, Got: zxvk}
		return
	}
	z.Network, err = dc.ReadString()
	if err != nil {
		return
	}
	z.SrcHost, err = dc.ReadString()
	if err != nil {
		return
	}
	z.SrcPort, err = dc.ReadUint16()
	if err != nil {
		return
	}
	z.DstHost, err = dc.ReadString()
	if err != nil {
		return
	}
	z.DstPort, err = dc.ReadUint16()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ConnInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 5
	err = en.Append(0x95)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Network)
	if err != nil {
		return
	}
	err = en.WriteString(z.SrcHost)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.SrcPort)
	if err != nil {
		return
	}
	err = en.WriteString(z.DstHost)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.DstPort)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ConnInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 5
	o = append(o, 0x95)
	o = msgp.AppendString(o, z.Network)
	o = msgp.AppendString(o, z.SrcHost)
	o = msgp.AppendUint16(o, z.SrcPort)
	o = msgp.AppendString(o, z.DstHost)
	o = msgp.AppendUint16(o, z.DstPort)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ConnInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zbzg uint32
	zbzg, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zbzg != 5 {
		err = msgp.ArrayError{Wanted: 5, Got: zbzg}
		return
	}
	z.Network, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	z.SrcHost, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	z.SrcPort, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	z.DstHost, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	z.DstPort, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ConnInfo) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Network) + msgp.StringPrefixSize + len(z.SrcHost) + msgp.Uint16Size + msgp.StringPrefixSize + len(z.DstHost) + msgp.Uint16Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameData) DecodeMsg(dc *msgp.Reader) (err error) {
	var zbai uint32
	zbai, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zbai != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zbai}
		return
	}
	z.StreamId, err = dc.ReadUint16()
	if err != nil {
		return
	}
	z.Data, err = dc.ReadBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameData) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return err
	}
	err = en.WriteUint16(z.StreamId)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendUint16(o, z.StreamId)
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zcmr uint32
	zcmr, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zcmr != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zcmr}
		return
	}
	z.StreamId, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameData) Msgsize() (s int) {
	s = 1 + msgp.Uint16Size + msgp.BytesPrefixSize + len(z.Data)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameDns) DecodeMsg(dc *msgp.Reader) (err error) {
	var zajw uint32
	zajw, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zajw != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zajw}
		return
	}
	z.Data, err = dc.ReadBytes(z.Data)
	if err != nil {
		return
	}
	z.StreamId, err = dc.ReadUint16()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameDns) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.StreamId)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameDns) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendBytes(o, z.Data)
	o = msgp.AppendUint16(o, z.StreamId)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameDns) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zwht uint32
	zwht, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zwht != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zwht}
		return
	}
	z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
	if err != nil {
		return
	}
	z.StreamId, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameDns) Msgsize() (s int) {
	s = 1 + msgp.BytesPrefixSize + len(z.Data) + msgp.Uint16Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameFin) DecodeMsg(dc *msgp.Reader) (err error) {
	var zhct uint32
	zhct, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zhct != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zhct}
		return
	}
	z.StreamId, err = dc.ReadUint16()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z FrameFin) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	err = en.WriteUint16(z.StreamId)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z FrameFin) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	o = msgp.AppendUint16(o, z.StreamId)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameFin) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zcua uint32
	zcua, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zcua != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zcua}
		return
	}
	z.StreamId, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z FrameFin) Msgsize() (s int) {
	s = 1 + msgp.Uint16Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameRst) DecodeMsg(dc *msgp.Reader) (err error) {
	var zxhx uint32
	zxhx, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zxhx != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zxhx}
		return
	}
	z.StreamId, err = dc.ReadUint16()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z FrameRst) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	err = en.WriteUint16(z.StreamId)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z FrameRst) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	o = msgp.AppendUint16(o, z.StreamId)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameRst) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zlqf uint32
	zlqf, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zlqf != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zlqf}
		return
	}
	z.StreamId, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z FrameRst) Msgsize() (s int) {
	s = 1 + msgp.Uint16Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameSyn) DecodeMsg(dc *msgp.Reader) (err error) {
	var zdaf uint32
	zdaf, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zdaf != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zdaf}
		return
	}
	z.StreamId, err = dc.ReadUint16()
	if err != nil {
		return
	}
	err = z.Address.DecodeMsg(dc)
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameSyn) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return err
	}
	err = en.WriteUint16(z.StreamId)
	if err != nil {
		return
	}
	err = z.Address.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameSyn) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendUint16(o, z.StreamId)
	o, err = z.Address.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameSyn) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zpks uint32
	zpks, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zpks != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zpks}
		return
	}
	z.StreamId, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	bts, err = z.Address.UnmarshalMsg(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameSyn) Msgsize() (s int) {
	s = 1 + msgp.Uint16Size + z.Address.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameSynResult) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zjfb uint32
	zjfb, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zjfb > 0 {
		zjfb--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "StreamId":
			z.StreamId, err = dc.ReadUint16()
			if err != nil {
				return
			}
		case "Errno":
			z.Errno, err = dc.ReadUint32()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z FrameSynResult) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "StreamId"
	err = en.Append(0x82, 0xa8, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteUint16(z.StreamId)
	if err != nil {
		return
	}
	// write "Errno"
	err = en.Append(0xa5, 0x45, 0x72, 0x72, 0x6e, 0x6f)
	if err != nil {
		return err
	}
	err = en.WriteUint32(z.Errno)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z FrameSynResult) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "StreamId"
	o = append(o, 0x82, 0xa8, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64)
	o = msgp.AppendUint16(o, z.StreamId)
	// string "Errno"
	o = append(o, 0xa5, 0x45, 0x72, 0x72, 0x6e, 0x6f)
	o = msgp.AppendUint32(o, z.Errno)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameSynResult) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcxo uint32
	zcxo, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcxo > 0 {
		zcxo--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "StreamId":
			z.StreamId, bts, err = msgp.ReadUint16Bytes(bts)
			if err != nil {
				return
			}
		case "Errno":
			z.Errno, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z FrameSynResult) Msgsize() (s int) {
	s = 1 + 9 + msgp.Uint16Size + 6 + msgp.Uint32Size
	return
}
