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

	"github.com/ezex-io/proxier/config"
	"github.com/ezex-io/proxier/internal/server"
	"github.com/ezex-io/proxier/version"
	_ "go.uber.org/automaxprocs"
)

func main() {
	log := slog.Default()

	configPath := flag.String("config", "./config.yaml", "Path to config file")
	ver := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *ver {
		fmt.Println("Proxier v" + version.Version.String())
		os.Exit(0)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	log.Info("configuration loaded successfully")

	srv, err := server.New(log, cfg.Server, cfg.Proxy)
	if err != nil {
		log.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	srv.Start()
	log.Info("server started", "address", cfg.Server.Host+":"+cfg.Server.ListenPort)

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
