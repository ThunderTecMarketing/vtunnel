package msocks

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *FrameBase) DecodeMsg(dc *msgp.Reader) (err error) {
	var zxvk uint32
	zxvk, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zxvk != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zxvk}
		return
	}
	z.Streamid, err = dc.ReadUint16()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z FrameBase) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	err = en.WriteUint16(z.Streamid)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z FrameBase) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	o = msgp.AppendUint16(o, z.Streamid)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameBase) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zbzg uint32
	zbzg, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zbzg != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zbzg}
		return
	}
	z.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z FrameBase) Msgsize() (s int) {
	s = 1 + msgp.Uint16Size
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
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zcmr uint32
		zcmr, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zcmr != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zcmr}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
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
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
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
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zajw uint32
	zajw, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zajw != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zajw}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zwht uint32
		zwht, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zwht != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zwht}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
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
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	s += msgp.BytesPrefixSize + len(z.Data)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameDns) DecodeMsg(dc *msgp.Reader) (err error) {
	var zhct uint32
	zhct, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zhct != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zhct}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zcua uint32
		zcua, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zcua != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zcua}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	z.Data, err = dc.ReadBytes(z.Data)
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
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	err = en.WriteBytes(z.Data)
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
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameDns) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zxhx uint32
	zxhx, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zxhx != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zxhx}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zlqf uint32
		zlqf, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zlqf != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zlqf}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameDns) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	s += msgp.BytesPrefixSize + len(z.Data)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameFin) DecodeMsg(dc *msgp.Reader) (err error) {
	var zdaf uint32
	zdaf, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zdaf != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zdaf}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zpks uint32
		zpks, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zpks != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zpks}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameFin) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameFin) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameFin) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zjfb uint32
	zjfb, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zjfb != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zjfb}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zcxo uint32
		zcxo, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zcxo != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zcxo}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameFin) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FramePing) DecodeMsg(dc *msgp.Reader) (err error) {
	var zeff uint32
	zeff, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zeff != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zeff}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zrsw uint32
		zrsw, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zrsw != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zrsw}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FramePing) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FramePing) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FramePing) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zxpk uint32
	zxpk, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zxpk != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zxpk}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zdnj uint32
		zdnj, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zdnj != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zdnj}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FramePing) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameSynResult) DecodeMsg(dc *msgp.Reader) (err error) {
	var zobc uint32
	zobc, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zobc != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zobc}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zsnv uint32
		zsnv, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zsnv != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zsnv}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	z.Errno, err = dc.ReadUint32()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameSynResult) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return err
	}
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	err = en.WriteUint32(z.Errno)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameSynResult) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	o = msgp.AppendUint32(o, z.Errno)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameSynResult) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zkgt uint32
	zkgt, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zkgt != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zkgt}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zema uint32
		zema, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zema != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zema}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	z.Errno, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameSynResult) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	s += msgp.Uint32Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameRst) DecodeMsg(dc *msgp.Reader) (err error) {
	var zpez uint32
	zpez, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zpez != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zpez}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zqke uint32
		zqke, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zqke != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zqke}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameRst) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameRst) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameRst) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zqyh uint32
	zqyh, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zqyh != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zqyh}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zyzr uint32
		zyzr, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zyzr != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zyzr}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameRst) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameSpam) DecodeMsg(dc *msgp.Reader) (err error) {
	var zywj uint32
	zywj, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zywj != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zywj}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zjpj uint32
		zjpj, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if zjpj != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zjpj}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	z.Data, err = dc.ReadBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameSpam) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return err
	}
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameSpam) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameSpam) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zzpf uint32
	zzpf, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zzpf != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zzpf}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zrfe uint32
		zrfe, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zrfe != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zrfe}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameSpam) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	s += msgp.BytesPrefixSize + len(z.Data)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FrameSyn) DecodeMsg(dc *msgp.Reader) (err error) {
	var zgmo uint32
	zgmo, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zgmo != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zgmo}
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var ztaf uint32
		ztaf, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if ztaf != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: ztaf}
			return
		}
		z.FrameBase.Streamid, err = dc.ReadUint16()
		if err != nil {
			return
		}
	}
	z.Network, err = dc.ReadString()
	if err != nil {
		return
	}
	z.Address, err = dc.ReadString()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FrameSyn) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 3
	err = en.Append(0x93)
	if err != nil {
		return err
	}
	if z.FrameBase == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// array header, size 1
		err = en.Append(0x91)
		if err != nil {
			return err
		}
		err = en.WriteUint16(z.FrameBase.Streamid)
		if err != nil {
			return
		}
	}
	err = en.WriteString(z.Network)
	if err != nil {
		return
	}
	err = en.WriteString(z.Address)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FrameSyn) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 3
	o = append(o, 0x93)
	if z.FrameBase == nil {
		o = msgp.AppendNil(o)
	} else {
		// array header, size 1
		o = append(o, 0x91)
		o = msgp.AppendUint16(o, z.FrameBase.Streamid)
	}
	o = msgp.AppendString(o, z.Network)
	o = msgp.AppendString(o, z.Address)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FrameSyn) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zeth uint32
	zeth, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zeth != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zeth}
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.FrameBase = nil
	} else {
		if z.FrameBase == nil {
			z.FrameBase = new(FrameBase)
		}
		var zsbz uint32
		zsbz, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			return
		}
		if zsbz != 1 {
			err = msgp.ArrayError{Wanted: 1, Got: zsbz}
			return
		}
		z.FrameBase.Streamid, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			return
		}
	}
	z.Network, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	z.Address, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FrameSyn) Msgsize() (s int) {
	s = 1
	if z.FrameBase == nil {
		s += msgp.NilSize
	} else {
		s += 1 + msgp.Uint16Size
	}
	s += msgp.StringPrefixSize + len(z.Network) + msgp.StringPrefixSize + len(z.Address)
	return
}
