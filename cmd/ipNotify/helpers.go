package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fasonju/ipNotify/internal/notify"
	"github.com/fasonju/ipNotify/internal/requests"
	"github.com/fasonju/ipNotify/internal/types"
)

func listenIps(cfg *types.Config) {
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

func getInitialIPs(cfg *types.Config) (string, string) {
	var previousIpv4, previousIpv6 string

	if cfg.Ipv4Enabled {
		ip, err := requests.GetIP(cfg.Ipv4url)
		if err != nil {
			slog.Error("Failed to get Ipv4", "error", err)
		} else {
			previousIpv4 = ip
			slog.Info("Initial Ipv4 queried", "ipv4", ip)
		}
	}

	if cfg.Ipv6Enabled {
		ip, err := requests.GetIP(cfg.Ipv6url)
		if err != nil {
			slog.Error("Failed to get Ipv6", "error", err)
		} else {
			previousIpv6 = ip
			slog.Info("Initial Ipv6 queried", "ipv6", ip)
		}
	}

	return previousIpv4, previousIpv6
}

func setupSignalChannel() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}

func checkIpDiffAndNotify(previousIpv4, previousIpv6 string, cfg *types.Config) (string, string, error) {
	newIpv4, newIpv6, err := fetchCurrentIPs(cfg)
	if err != nil {
		return previousIpv4, previousIpv6, err
	}

	if ipsChanged(previousIpv4, newIpv4, previousIpv6, newIpv6) {
		message := formatChangeMessage(previousIpv4, newIpv4, previousIpv6, newIpv6, cfg)
		if cfg.SmtpEnabled {
			slog.Info("Notifying through SMTP")
			notify.NotifySMTP(cfg, message)
		}
	} else {
		slog.Info("No IP change", "ipv4", newIpv4, "ipv6", newIpv6)
	}

	return newIpv4, newIpv6, nil
}

func fetchCurrentIPs(cfg *types.Config) (string, string, error) {
	var newIpv4, newIpv6 string
	var err error

	if cfg.Ipv4Enabled {
		newIpv4, err = requests.GetIP(cfg.Ipv4url)
		if err != nil {
			return "", "", err
		}
	}

	if cfg.Ipv6Enabled {
		newIpv6, err = requests.GetIP(cfg.Ipv6url)
		if err != nil {
			return "", "", err
		}
	}

	return newIpv4, newIpv6, nil
}

func ipsChanged(prev4, new4, prev6, new6 string) bool {
	return new4 != prev4 || new6 != prev6
}

func formatChangeMessage(prev4, new4, prev6, new6 string, cfg *types.Config) string {
	var builder strings.Builder

	if cfg.Ipv4Enabled && new4 != prev4 {
		builder.WriteString(fmt.Sprintf("IPv4 changed from %s to %s\n", prev4, new4))
	}
	if cfg.Ipv6Enabled && new6 != prev6 {
		builder.WriteString(fmt.Sprintf("IPv6 changed from %s to %s\n", prev6, new6))
	}

	return builder.String()
}
