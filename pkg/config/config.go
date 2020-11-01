package config

import (
	"context"
	"fmt"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

// Peer represents one connection from client to server
type Peer struct {
	NIC    string `json:"nic"`
	Listen string `json:"listen"`
	Server string `json:"server"`
}

// ClientConfig - For client
type ClientConfig struct {
	HTTPProxy string `json:"http_proxy"`

	Peers []Peer `json:"peers"`
}

// ServerConfig - For server
type ServerConfig struct {
	Endpoints []string `json:"endpoints"`
}

const (
	// ModeServer - server in Config.Mode
	ModeServer = "server"

	// ModeClient - client in Config.Mode
	ModeClient = "client"
)

// Config parses all configs from files
type Config struct {
	Client ClientConfig `json:"client"`
	Server ServerConfig `json:"server"`

	Mode string `json:"mode"`
}

// ReadConfig loads config from files
func ReadConfig(ctx context.Context) (*Config, error) {
	loader := confita.NewLoader(
		file.NewOptionalBackend("/etc/mpproxy/mpproxy.json"),
		file.NewOptionalBackend("/etc/mpproxy/mpproxy.yaml"),
		file.NewOptionalBackend("./mpproxy.json"),
		file.NewOptionalBackend("./mpproxy.yaml"),
	)

	cfg := &Config{}

	if err := loader.Load(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("failed to laod config: %w", err)
	}

	return cfg, nil
}
