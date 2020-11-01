package runner

import (
	"context"
	"time"

	"github.com/cs3238-tsuzu/multipath-proxy/pkg/config"
)

// Run loads config and starts client/server
func Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.ReadConfig(ctx)

	if err != nil {
		panic(err)
	}

	switch cfg.Mode {
	case config.ModeServer:
		err = runServer(cfg)
	case config.ModeClient:
		err = runClient(cfg)
	default:
		panic("unknown mode: " + cfg.Mode)
	}

	if err != nil {
		panic(err)
	}
}
