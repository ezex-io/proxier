package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ezex-io/gopkg/logger"
	"github.com/ezex-io/proxier/internal/proxy"
)

type httpServer struct {
	httpServer *http.Server
	errCh      chan error
	log        logger.Logger
}

func NewHTTP(log logger.Logger, address string, proxyRules map[string]string) (Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Proxier is running"))
	})

	mux.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	for endpoint, destination := range proxyRules {
		router, handler, err := proxy.HTTPHandler(endpoint, destination)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxy handler for endpoint %s: %w", endpoint, err)
		}

		mux.Handle(router+"/", handler) // Ensure `/route/` works
		mux.Handle(router, handler)     // Ensure `/route` works
		log.Info("Registered proxy route", "endpoint", endpoint, "destination", destination)
	}

	srv := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // Prevent slow client attacks
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second, // Keep connections alive for long-lived clients
	}

	return &httpServer{
		httpServer: srv,
		errCh:      make(chan error, 1),
		log:        log,
	}, nil
}

func (s *httpServer) Start() {
	go func() {
		s.log.Info("starting server", "address", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errCh <- fmt.Errorf("server error: %w", err)
		}
	}()
}

func (s *httpServer) Notify() <-chan error {
	return s.errCh
}

func (s *httpServer) Stop(ctx context.Context) {
	s.log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.log.Error("server forced to shutdown", "error", err)
	} else {
		s.log.Info("server gracefully stopped")
	}
}
