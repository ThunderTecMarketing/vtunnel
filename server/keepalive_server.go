/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package server

import (
	"time"
	"os"
	"net"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// Borrowed from the Go standard library.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts the connection with a keep-alive enabled.
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// File implements caddy.Listener; it returns the underlying file of the listener.
func (ln tcpKeepAliveListener) File() (*os.File, error) {
	return ln.TCPListener.File()
}
