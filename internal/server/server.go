package server

import (
	"context"
	"log/slog"

	"github.com/ezex-io/proxier/config"
)

type Server interface {
	Start()
	Notify() <-chan error
	Stop(ctx context.Context)
}

func New(cfg *config.Config, log *slog.Logger) (Server, error) {
	if cfg.Server.FastHTTP {
		return newFastHTTP(log, cfg.Server, cfg.Proxy)
	}

	return NewHTTP(log, cfg.Server, cfg.Proxy)
}
