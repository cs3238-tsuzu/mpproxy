package server

import (
	"crypto/tls"
	"fmt"
	"net"

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

	tlsConf, err := generateTLSConfig()
	config := (*quic.Config)(nil)

	if err != nil {
		return nil, fmt.Errorf("failed to generate TLS config: %w", err)
	}

	listeners := make([]net.Listener, len(endpoints))
	for i := range endpoints {
		listener, err := server.listen(endpoints[i], tlsConf, config)

		if err != nil {
			return nil, fmt.Errorf("failed to listen on %s: %w", endpoints[i], err)
		}

		listeners = append(listeners, netutil.NewFlattenListener(listener))
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
