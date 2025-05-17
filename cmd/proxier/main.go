package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ezex-io/gopkg/env"
	"github.com/ezex-io/gopkg/logger"
	"github.com/ezex-io/proxier/config"
	"github.com/ezex-io/proxier/internal/server"
	"github.com/ezex-io/proxier/version"
	_ "go.uber.org/automaxprocs"
)

func main() {
	log := logger.NewSlog(logger.WithTextHandler(os.Stdout, slog.LevelDebug))

	envFile := flag.String("env", ".env", "Path to environment file")
	ver := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *ver {
		fmt.Println("Proxier v" + version.Version.String())
		os.Exit(0)
	}

	if err := env.LoadEnvsFromFile(*envFile); err != nil {
		log.Debug("Failed to load env file '%s': %v. Continuing with system environment...", *envFile, err)
	}

	cfg := config.LoadFromEnv()

	log.Info("configuration loaded successfully")

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	srv.Start()
	log.Info("server started", "address", cfg.Address)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		log.Warn("termination signal received", "signal", sig.String())

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		srv.Stop(shutdownCtx)
	case err := <-srv.Notify():
		log.Error("server encountered an error", "error", err)
	}
}
