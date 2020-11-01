package client

import (
	"context"
	"net"

	"github.com/getlantern/multipath"
	"github.com/lucas-clemente/quic-go"
)

type quicStreamConn struct {
	label   string
	session quic.Session
	quic.Stream
}

func newQuicStreamConn(
	label string,
	session quic.Session,
	stream quic.Stream,
) *quicStreamConn {
	return &quicStreamConn{
		label:   label,
		session: session,
		Stream:  stream,
	}
}

var _ multipath.Dialer = &quicStreamConn{}
var _ net.Conn = &quicStreamConn{}

func (qsc *quicStreamConn) DialContext(ctx context.Context) (net.Conn, error) {
	return qsc, nil
}

func (qsc *quicStreamConn) Label() string {
	return qsc.label
}

// LocalAddr returns the local network address.
func (qsc *quicStreamConn) LocalAddr() net.Addr {
	return qsc.session.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (qsc *quicStreamConn) RemoteAddr() net.Addr {
	return qsc.session.RemoteAddr()
}
