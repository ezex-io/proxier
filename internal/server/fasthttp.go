package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ezex-io/proxier/config"
	"github.com/ezex-io/proxier/internal/proxy"
	"github.com/valyala/fasthttp"
)

type fastHTTPServer struct {
	sv     *fasthttp.Server
	errCh  chan error
	log    *slog.Logger
	addr   string
	cancel context.CancelFunc
}

func newFastHTTP(log *slog.Logger, cfg *config.ServerConfig, proxyRules []*config.ProxyRule) (Server, error) {
	handlers := make(map[string]fasthttp.RequestHandler)

	for _, rule := range proxyRules {
		endpoint, handler, err := proxy.FastHTTPHandler(rule.Endpoint, rule.DestinationURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create fasthttp proxy handler for %s: %w", rule.Endpoint, err)
		}
		handlers[endpoint] = handler
		log.Info("Registered proxy route", "endpoint", endpoint, "destination", rule.DestinationURL)
	}

	handler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		switch path {
		case "/":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString("Proxier is running")

			return
		case "/livez":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString("OK")

			return
		}

		for endpoint, h := range handlers {
			if strings.HasPrefix(path, endpoint) {
				h(ctx)

				return
			}
		}

		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString("Route not found")
	}

	srv := &fastHTTPServer{
		sv: &fasthttp.Server{
			Handler: handler,

			// Optimized settings
			Name:                  "proxier-fasthttp",
			Concurrency:           0,
			ReadBufferSize:        4096,
			WriteBufferSize:       4096,
			ReadTimeout:           10 * time.Second,
			WriteTimeout:          15 * time.Second,
			IdleTimeout:           60 * time.Second,
			MaxRequestsPerConn:    1000,
			MaxConnsPerIP:         100,
			MaxRequestBodySize:    10 * 1024 * 1024,
			MaxIdleWorkerDuration: 10 * time.Second,
			TCPKeepalive:          true,
			ReduceMemoryUsage:     true,
			DisableKeepalive:      false,
			StreamRequestBody:     true,
			LogAllErrors:          false,
			SecureErrorLogMessage: true,
			ErrorHandler: func(ctx *fasthttp.RequestCtx, _ error) {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetBodyString("Internal Server Error")
			},
		},
		errCh: make(chan error, 1),
		log:   log,
		addr:  fmt.Sprintf("%s:%s", cfg.Host, cfg.ListenPort),
	}

	return srv, nil
}

func (s *fastHTTPServer) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	go func() {
		s.log.Info("starting fasthttp server", "address", s.addr)
		if err := s.sv.ListenAndServe(s.addr); err != nil {
			s.errCh <- fmt.Errorf("fasthttp server error: %w", err)
		}
		<-ctx.Done()
	}()
}

func (s *fastHTTPServer) Notify() <-chan error {
	return s.errCh
}

func (s *fastHTTPServer) Stop(_ context.Context) {
	s.log.Info("shutting down fasthttp server...")

	s.cancel()
	if err := s.sv.Shutdown(); err != nil {
		s.log.Error("failed to shutdown fasthttp server", "error", err)
	} else {
		s.log.Info("fasthttp server stopped")
	}
}
