package main

import (
	"log/slog"
	"os"

	"github.com/fasonju/ipNotify/internal/config"
)

func main() {
	slog.Info("ipNotify started, loading config")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("config loaded")
	listenIps(cfg)
}
