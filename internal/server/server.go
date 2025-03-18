package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ezex-io/proxier/config"
	"github.com/ezex-io/proxier/internal/proxy"
)

type Server struct {
	httpServer *http.Server
	errCh      chan error
	log        *slog.Logger
}

func New(log *slog.Logger, serverCfg *config.ServerConfig, proxyRules []*config.ProxyRule) (*Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Proxier is running"))
	})

	mux.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	for _, rule := range proxyRules {
		endpoint, handler, err := proxy.NewProxyHandler(rule.Endpoint, rule.DestinationURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxy handler for endpoint %s: %w", endpoint, err)
		}

		mux.Handle(endpoint+"/", handler) // Ensure `/route/` works
		mux.Handle(endpoint, handler)     // Ensure `/route` works
		log.Info("Registered proxy route", "endpoint", rule.Endpoint, "destination", rule.DestinationURL)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", serverCfg.Host, serverCfg.ListenPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // Prevent slow client attacks
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second, // Keep connections alive for long-lived clients
	}

	return &Server{
		httpServer: srv,
		errCh:      make(chan error, 1),
		log:        log,
	}, nil
}

func (s *Server) Start() {
	go func() {
		s.log.Info("starting server", "address", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errCh <- fmt.Errorf("server error: %w", err)
		}
	}()
}

func (s *Server) Notify() <-chan error {
	return s.errCh
}

func (s *Server) HTTPHandler() http.Handler {
	return s.httpServer.Handler
}

func (s *Server) Stop(ctx context.Context) {
	s.log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.log.Error("server forced to shutdown", "error", err)
	} else {
		s.log.Info("server gracefully stopped")
	}
}
