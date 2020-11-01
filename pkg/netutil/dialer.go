package netutil

import (
	"context"
	"net"

	"github.com/getlantern/multipath"
)

// NewDialer returns multipath.Dialer wrapping existing net.Conn
func NewDialer(label string, conn net.Conn) multipath.Dialer {
	return &dialer{
		label: label,
		conn:  conn,
	}
}

type dialer struct {
	label string
	conn  net.Conn
}

var _ multipath.Dialer = &dialer{}

// DialContext connects to the server.
func (d *dialer) DialContext(ctx context.Context) (net.Conn, error) {
	return d.conn, nil
}

// Label returns a label for multipath.Dialer
func (d *dialer) Label() string {
	return d.label
}
