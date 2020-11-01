package config

import (
	"context"
	"fmt"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

// Peer represents one connection from client to server
type Peer struct {
	NIC    string `json:"nic" yaml:"nic" config:"nic"`
	Listen string `json:"listen" yaml:"listen" config:"listen"`
	Server string `json:"server" yaml:"server" config:"server"`
}

// ClientConfig - For client
type ClientConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint" config:"endpoints"`

	Peers []Peer `json:"peers" yaml:"peers" config:"peers"`
}

// ServerConfig - For server
type ServerConfig struct {
	Backend   string   `json:"backend" yaml:"backend" config:"backend"`
	Endpoints []string `json:"endpoints" yaml:"endpoints" config:"endpoints"`
}

const (
	// ModeServer - server in Config.Mode
	ModeServer = "server"

	// ModeClient - client in Config.Mode
	ModeClient = "client"
)

// Config parses all configs from files
type Config struct {
	Client ClientConfig `json:"client" yaml:"client"`
	Server ServerConfig `json:"server" yaml:"server"`

	Mode string `json:"mode" yaml:"mode" config:"mode"`
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

	if err := loader.Load(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to laod config: %w", err)
	}

	return cfg, nil
}
