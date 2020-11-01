package netutil

import (
	"net"

	"github.com/lucas-clemente/quic-go"
)

// NewStreamConn returns a new net.Conn, which is a wrapper for QUIC stream/session as net.Conn
func NewStreamConn(
	session quic.Session,
	stream quic.Stream,
) net.Conn {
	return &streamConn{
		session: session,
		Stream:  stream,
	}
}

type streamConn struct {
	session quic.Session
	quic.Stream
}

var _ net.Conn = &streamConn{}

// LocalAddr returns the local network address.
func (qsc *streamConn) LocalAddr() net.Addr {
	return qsc.session.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (qsc *streamConn) RemoteAddr() net.Addr {
	return qsc.session.RemoteAddr()
}
