/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package vpn

import (
	"encoding/binary"
	"io"
)

var DefaultMTU = 1400

type PacketType uint8

const PacketTypeDNS = 1
const PacketTypeOPEN = 2
const PacketTypePACKETS = 3

type Packet interface {
}

type BasicPacket struct {
	packetType PacketType
	data       []byte
}

func (p *BasicPacket) Unpack(r io.Reader) (err error) {

	buf0 := make([]byte, 1)
	buf1 := make([]byte, 2)
	buf2 := make([]byte, DefaultMTU)

	if _, err = io.ReadFull(r, buf0); err != nil {
		//maybe io.ErrUnexpectedEOF
		return
	}
	p.packetType = PacketType(buf0[0])

	if _, err = io.ReadFull(r, buf1); err != nil {
		//maybe io.ErrUnexpectedEOF
		return
	}

	packetLength := binary.BigEndian.Uint16(buf1)

	if int(packetLength) > DefaultMTU || packetLength <= 0 {
		return ErrPacketLengthInvalid
	}

	if _, err = io.ReadFull(r, buf2[:packetLength]); err != nil {
		return
	}

	return nil
}

func (p *BasicPacket) Pack() []byte {

	lBuf := make([]byte, 3 + len(p.data))

	lBuf[0] = byte(p.packetType)
	binary.BigEndian.PutUint16(lBuf[1:3], uint16(len(p.data)))
	copy(lBuf[3:], p.data)

	return nil
}