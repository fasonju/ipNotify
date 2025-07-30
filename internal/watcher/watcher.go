package watcher

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fasonju/ipNotify/internal/types"
)

// ListenIps starts a loop that periodically checks for IP changes and handles graceful shutdown.
func ListenIps(cfg *types.Config) {
	previousIpv4, previousIpv6 := getInitialIPs(cfg)

	sigs := setupSignalChannel()
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	slog.Info("Starting ip watcher loop", "interval", cfg.Interval)

	for {
		select {
		case t := <-ticker.C:
			slog.Info("Checking for IP diffs", "time", t)
			var err error
			previousIpv4, previousIpv6, err = checkIpDiffAndNotify(previousIpv4, previousIpv6, cfg)
			if err != nil {
				slog.Error("Unable to check IPs", "error", err)
			}
		case sig := <-sigs:
			slog.Info("Received signal", "signal", sig.String())
			slog.Info("Shutting down gracefully")
			return
		}
	}
}

// setupSignalChannel sets up a channel to listen for OS interrupt and terminate signals.
func setupSignalChannel() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}
