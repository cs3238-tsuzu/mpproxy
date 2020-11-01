package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/cs3238-tsuzu/multipath-proxy/pkg/config"
	"github.com/cs3238-tsuzu/multipath-proxy/pkg/netutil"
	"github.com/cs3238-tsuzu/multipath-proxy/pkg/nic"
	"github.com/getlantern/multipath"
	"github.com/lucas-clemente/quic-go"
)

// Client establishes multipath connections to the destination server
type Client struct {
	sessions []quic.Session
}

// NewClient returns a new client
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	if cfg.Mode != config.ModeClient {
		return nil, fmt.Errorf("config mode should be client")
	}

	peers, err := nic.GetPeers(cfg.Client.Peers)

	if err != nil {
		return nil, fmt.Errorf("failed to get peers: %w", err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}

	config := (*quic.Config)(nil)

	client := &Client{
		sessions: make([]quic.Session, 0, len(peers)),
	}

	for _, peer := range peers {
		session, err := client.conncet(ctx, peer, tlsConf, config)

		if err != nil {
			return nil, fmt.Errorf("failed to connec to %s: %w", peer.Server, err)
		}

		client.sessions = append(client.sessions, session)
	}

	return client, nil
}

func (c *Client) conncet(ctx context.Context, peer *nic.Peer, tlsConf *tls.Config, config *quic.Config) (quic.Session, error) {
	pconn, err := net.ListenUDP("", peer.Address)

	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s for %s: %w", peer.Address, peer.Server, err)
	}

	addr, err := net.ResolveIPAddr("", peer.Server)

	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address(%s): %w", peer.Server, err)
	}

	session, err := quic.DialContext(ctx, pconn, addr, peer.Server, tlsConf, config)

	if err != nil {
		return nil, fmt.Errorf("failed to establish quic connection: %w", err)
	}

	return session, nil
}

// NewMultipathConn returns a new multipath client powered by QUIC stream
func (c *Client) NewMultipathConn(ctx context.Context) (net.Conn, error) {
	dialers, err := c.newDialers()

	if err != nil {
		return nil, fmt.Errorf("failed to initialize errors: %w", err)
	}

	dialer := multipath.NewDialer("mpproxy-server", dialers)

	conn, err := dialer.DialContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize multipath connection: %w", err)
	}

	return conn, nil
}

func (c *Client) newDialers() ([]multipath.Dialer, error) {
	dialers := make([]multipath.Dialer, len(c.sessions))

	for i := range c.sessions {
		stream, err := c.sessions[i].OpenStream()

		if err != nil {
			return nil, fmt.Errorf("failed to open stream for %s: %w",
				c.sessions[i].RemoteAddr().String(),
				err,
			)
		}

		dialers = append(dialers,
			netutil.NewDialer(
				fmt.Sprintf("%d", i), // TODO: Use another label
				netutil.NewStreamConn(
					c.sessions[i],
					stream,
				),
			),
		)
	}

	return dialers, nil
}