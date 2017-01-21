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
	"github.com/FTwOoO/netstack/tcpip/buffer"
	"io"
	"encoding/binary"
)

var DefaultMTU = 1400

func ReadPackets(r io.Reader) (packets []buffer.View, err error) {

	buf1 := make([]byte, 2)
	buf2 := make([]byte, DefaultMTU)

	for {
		if _, err = io.ReadFull(r, buf1); err != nil {
			//nothing to read
			if err == io.EOF {
				err = nil
				break
			}
			//maybe io.ErrUnexpectedEOF
			return nil, err
		}

		packetLength := binary.BigEndian.Uint16(buf1)

		if int(packetLength) > DefaultMTU || packetLength <= 0 {
			return nil, ErrPacketLengthInvalid
		}

		if _, err = io.ReadFull(r, buf2[:packetLength]); err != nil {
			return nil, err
		}

		newPacket := make([]byte, packetLength)
		copy(newPacket, buf2[:packetLength])
		packets = append(packets, newPacket)
	}

	if len(packets) <= 0 {
		return nil, io.ErrUnexpectedEOF
	}

	return
}

func WritePackets(w io.Writer, packets []buffer.View) error {
	lBuf := make([]byte, 2)

	for _, packet := range packets {
		binary.BigEndian.PutUint16(lBuf, uint16(len(packet)))
		n, err := w.Write(lBuf)
		if n != 2 && err == nil {
			err = ErrWriteFail
		}

		if err != nil {
			return err
		}

		n, err = w.Write([]byte(packet))
		if n != len(packet) && err == nil {
			err = ErrWriteFail
		}

		if err != nil {
			return err
		}
	}

	return nil
}