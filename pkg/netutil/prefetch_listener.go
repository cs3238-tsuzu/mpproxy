package netutil

import (
	"bytes"
	"io"
	"net"
	"time"
)

type connError struct {
	conn net.Conn
	err  error
}

// prefetchListener returns new connections only if any data are received
type prefetchListener struct {
	net.Listener
	expected int
	timeout  time.Duration

	endCh  chan struct{}
	connCh chan connError
}

// NewPrefetchListener returns a listener that returns new connections only if any data are received
func NewPrefetchListener(l net.Listener, expected int, timeout time.Duration) net.Listener {
	pl := &prefetchListener{
		endCh:  make(chan struct{}, 1),
		connCh: make(chan connError, 0),

		Listener: l,
		expected: expected,
		timeout:  timeout,
	}

	pl.start()

	return pl
}

func (pl *prefetchListener) putConnCh(conn net.Conn, err error) bool {
	timer := time.NewTimer(pl.timeout)
	defer timer.Stop()

	select {
	case pl.connCh <- connError{
		conn,
		nil,
	}:
		return true

	case <-timer.C:
	case <-pl.endCh:
	}

	return false
}

func (pl *prefetchListener) interrupted() bool {
	select {
	case <-pl.endCh:
		return true
	default:
	}

	return false
}

func (pl *prefetchListener) start() {
	go func() {
		for {
			if pl.interrupted() {
				return
			}

			conn, err := pl.Listener.Accept()

			if err != nil {
				if err, ok := err.(net.Error); ok {
					if !err.Temporary() {
						return
					}
				}

				pl.putConnCh(nil, err)

				continue
			}

			if pl.interrupted() {
				conn.Close()

				return
			}

			go func() {
				conn.SetReadDeadline(time.Now().Add(pl.timeout))
				buf := make([]byte, pl.expected)

				_, err := io.ReadFull(conn, buf)

				if err != nil || pl.interrupted() {
					pl.putConnCh(nil, err)
					conn.Close()

					return
				}

				conn.SetReadDeadline(time.Time{})

				ok := pl.putConnCh(
					&multiReadConn{conn, io.MultiReader(
						bytes.NewReader(buf),
						conn,
					)},
					nil,
				)

				if !ok {
					conn.Close()
				}
			}()
		}
	}()
}

func (pl *prefetchListener) Accept() (net.Conn, error) {
	if pl.interrupted() {
		return nil, ErrClosed
	}

	select {
	case <-pl.endCh:
		return nil, ErrClosed
	case c := <-pl.connCh:
		return c.conn, c.err
	}
}

func (pl *prefetchListener) Close() error {
	close(pl.endCh)

	return pl.Listener.Close()

}

type multiReadConn struct {
	net.Conn
	wrapped io.Reader
}

func (mrc *multiReadConn) Read(p []byte) (n int, err error) {
	return mrc.wrapped.Read(p)
}
