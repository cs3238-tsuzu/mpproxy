package runner

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/cs3238-tsuzu/multipath-proxy/pkg/config"
	"github.com/cs3238-tsuzu/multipath-proxy/pkg/server"
)

// runServer runs server-side listener
func runServer(cfg *config.Config) error {
	if cfg.Mode != config.ModeServer {
		return fmt.Errorf("config mode should be server")
	}

	listener, err := server.NewServer(cfg.Server.Endpoints)

	if err != nil {
		return fmt.Errorf("failed to initialize multipath client: %w", err)
	}

	for {
		mpconn, err := listener.Accept()

		if err != nil {
			if err, ok := err.(net.Error); ok {
				if !err.Temporary() {
					return fmt.Errorf("Accept returned a critical error: %w", err)
				}
			}

			return fmt.Errorf("failed to accept new multipath connection: %w", err)
		}

		go func() {
			conn, err := net.Dial("tcp", cfg.Server.Backend)

			if err != nil {
				log.Printf("failed to prepare connection to backend: %+v", err)

				return
			}

			go io.Copy(conn, mpconn)
			go io.Copy(mpconn, conn)
		}()
	}

}
