package server

import (
	"context"
	"strings"

	"github.com/ezex-io/gopkg/logger"
	"github.com/ezex-io/proxier/config"
)

type Server interface {
	Start()
	Notify() <-chan error
	Stop(ctx context.Context)
}

func New(cfg *config.Config, log logger.Logger) (Server, error) {
	rules := make(map[string]string)
	for _, rule := range cfg.ProxyRules {
		parsedRule := strings.Split(rule, "|")
		rules[parsedRule[0]] = parsedRule[1]
	}

	if cfg.EnableFastHTTP {
		return newFastHTTP(log, cfg.Address, rules)
	}

	return NewHTTP(log, cfg.Address, rules)
}
