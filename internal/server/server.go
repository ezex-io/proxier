package server

import (
	"context"

	"github.com/ezex-io/gopkg/logger"
	"github.com/ezex-io/proxier/config"
)

type Server interface {
	Start()
	Notify() <-chan error
	Stop(ctx context.Context)
}

func New(cfg *config.Config, log logger.Logger) (Server, error) {
	if cfg.EnableFastHTTP {
		return newFastHTTP(log, cfg.Address, cfg.ParsedRules)
	}

	return NewHTTP(log, cfg.Address, cfg.ParsedRules)
}
