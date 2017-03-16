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
	"github.com/FTwOoO/noise"
)

type NoiseIXHandshake struct {
	Hs            *noise.HandshakeState
	IsInitiator   bool
	Pattern       noise.HandshakePattern

	csEnc         *noise.CipherState
	csDec         *noise.CipherState
	handshakeStep int
}

func NewNoiseIXHandshake(cs noise.CipherSuite, pg []byte, staticKey noise.DHKey, isInitiator bool) (*NoiseIXHandshake, error) {

	this := new(NoiseIXHandshake)
	if cs == nil {
		cs = noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
	}

	this.Hs = noise.NewHandshakeState(noise.Config{
		CipherSuite: cs,
		Pattern: noise.HandshakeIX,
		Initiator: isInitiator,
		Prologue: pg,
		StaticKeypair: staticKey},
	)

	this.IsInitiator = isInitiator
	this.Pattern = noise.HandshakeIK
	return this, nil
}

func (this *NoiseIXHandshake) isHandshakeCompleted() bool {
	return this.handshakeStep >= len(this.Pattern.Messages)
}

func (this *NoiseIXHandshake) isWriteStep() bool {
	return (this.handshakeStep % 2 == 1 && this.IsInitiator) || (this.handshakeStep % 2 == 0 && !this.IsInitiator)
}

func (this *NoiseIXHandshake) isReadStep() bool {
	return (this.handshakeStep % 2 == 1 && !this.IsInitiator) || (this.handshakeStep % 2 == 0 && this.IsInitiator)
}

func (this *NoiseIXHandshake) isFinalStep() bool {
	return this.handshakeStep == len(this.Pattern.Messages)
}

func (this *NoiseIXHandshake) stepOne() {
	if !this.isHandshakeCompleted() {
		this.handshakeStep += 1
	}
}

func (this *NoiseIXHandshake) handshakeDone(csenc *noise.CipherState, csdec *noise.CipherState) {
	this.csEnc = csenc
	this.csDec = csdec
}

func (this *NoiseIXHandshake) Encode(b []byte) (packet []byte, err error) {

	if !this.isHandshakeCompleted() {
		this.stepOne()

		if this.isWriteStep() {
			var cs0, cs1 *noise.CipherState
			packet, cs0, cs1 = this.Hs.WriteMessage(nil, b)

			// final msg for noise handshake pattern we can extract the
			// cipher
			if this.isFinalStep() {
				this.handshakeDone(cs1, cs0)
			}

		} else {
			return nil, ErrInValidHandshakeStep
		}

	} else {
		packet = this.csEnc.EncryptWithAd(nil, nil, b)
	}

	return
}

func (this *NoiseIXHandshake) Decode(b []byte) (packet []byte, err error) {
	if !this.isHandshakeCompleted() {
		this.stepOne()

		if this.isReadStep() {
			var cs0, cs1 *noise.CipherState
			packet, cs0, cs1, err = this.Hs.ReadMessage(nil, b)
			if err != nil {
				return
			}

			// final msg for noise handshake pattern we can extract the
			// cipher
			if this.isFinalStep() {
				this.handshakeDone(cs0, cs1)
			}

			return

		} else {
			return nil, ErrInValidHandshakeStep
		}

	} else {
		packet, err = this.csDec.DecryptWithAd(nil, nil, b)
		if err != nil {
			return
		}

	}
	return
}
