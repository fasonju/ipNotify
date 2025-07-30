package watcher

import (
	"log/slog"

	"github.com/fasonju/ipNotify/internal/requests"
	"github.com/fasonju/ipNotify/internal/types"
)

// getInitialIPs fetches and logs the initial IP addresses (IPv4/IPv6) based on config.
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

// fetchCurrentIPs returns the current IPv4/IPv6 addresses or an error if a lookup fails.
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
