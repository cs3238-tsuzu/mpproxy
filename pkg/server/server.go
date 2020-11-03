package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/cs3238-tsuzu/multipath-proxy/pkg/netutil"
	"github.com/getlantern/multipath"
	"github.com/lucas-clemente/quic-go"
)

// Server handles multipath connections from client
type Server struct {
	net.Listener
}

var _ net.Listener = &Server{}

// NewServer inittializes a new server
func NewServer(endpoints []string) (*Server, error) {
	server := &Server{}

	config := &quic.Config{
		KeepAlive:          true,
		MaxIncomingStreams: 1 << 60,
	}

	tlsConf, err := generateTLSConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to generate TLS config: %w", err)
	}

	listeners := make([]net.Listener, len(endpoints))
	for i := range endpoints {
		listener, err := server.listen(endpoints[i], tlsConf, config)

		if err != nil {
			return nil, fmt.Errorf("failed to listen on %s: %w", endpoints[i], err)
		}

		listeners[i] = netutil.NewPrefetchListener(
			netutil.NewFlattenListener(listener),
			17,
			10*time.Second,
		)
	}

	stats := make([]multipath.StatsTracker, len(listeners))
	for i := range stats {
		stats[i] = multipath.NullTracker{}
	}

	server.Listener = multipath.NewListener(listeners, stats)

	return server, nil
}

func (s *Server) listen(endpoint string, tlsConf *tls.Config, config *quic.Config) (quic.Listener, error) {
	listener, err := quic.ListenAddr(endpoint, tlsConf, config)

	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", endpoint, err)
	}

	return listener, nil
}
