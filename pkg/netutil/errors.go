package netutil

import "fmt"

var (
	// ErrClosed represents the connection is already closed
	ErrClosed = fmt.Errorf("closed connection")
)
