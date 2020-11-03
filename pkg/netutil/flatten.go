package netutil

import (
	"context"
	"fmt"
	"net"

	"github.com/lucas-clemente/quic-go"
)

type flattenListener struct {
	listener quic.Listener
	endCh    chan struct{}
	accepted chan net.Conn
	errors   chan error
}

// NewFlattenListener returns a net.Listener for flattened QUIC listener
func NewFlattenListener(listener quic.Listener) net.Listener {
	l := &flattenListener{
		listener: listener,
		endCh:    make(chan struct{}, 1),
		errors:   make(chan error, 20),
		accepted: make(chan net.Conn, 1024),
	}

	l.start()

	return l
}

var _ net.Listener = &flattenListener{}

func (fl *flattenListener) saveError(err error) {
	select {
	case fl.errors <- err:
	default:
	}
}

func (fl *flattenListener) start() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-fl.endCh
		cancel()
	}()

	go func() {
		for {
			select {
			case <-fl.endCh:
				return
			default:
			}

			session, err := fl.listener.Accept(ctx)

			if err != nil {
				fl.saveError(fmt.Errorf("failed to accept QUIC session: %w", err))

				continue
			}

			go fl.runSession(ctx, session)
		}
	}()
}

func (fl *flattenListener) runSession(ctx context.Context, session quic.Session) {
	for {
		stream, err := session.AcceptStream(ctx)

		if err != nil {
			if err, ok := err.(net.Error); ok {
				if !err.Temporary() {
					return
				}
			}

			fl.saveError(fmt.Errorf("failed to accept QUIC stream: %w", err))

			continue
		}

		conn := NewStreamConn(session, stream)

		select {
		case fl.accepted <- conn:
		default:
			fl.saveError(fmt.Errorf("accepted channel is full"))

			conn.Close()
		}
	}
}

// Accept waits for and returns the next connection to the listener.
func (fl *flattenListener) Accept() (net.Conn, error) {
	select {
	case <-fl.endCh:

		return nil, ErrClosed
	default:
	}

	a := <-fl.accepted

	return a, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (fl *flattenListener) Close() error {
	close(fl.endCh)

	return nil
}

// Addr returns the listener's network address.
func (fl *flattenListener) Addr() net.Addr {
	return fl.listener.Addr()
}

func (fl *flattenListener) Errors() <-chan error {
	return fl.errors
}
