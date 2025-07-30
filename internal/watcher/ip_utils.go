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

// fetchCurrentIPs returns the current IPv4 and IPv6 addresses based on the config.
// It returns separate errors for IPv4 and IPv6 lookups if they fail.
// If an IP version is disabled, its address will be empty and error nil.
func fetchCurrentIPs(cfg *types.Config) (string, string, error, error) {
	var newIpv4, newIpv6 string
	var ipv4Err error
	var ipv6Err error

	if cfg.Ipv4Enabled {
		newIpv4, ipv4Err = requests.GetIP(cfg.Ipv4url)
	}

	if cfg.Ipv6Enabled {
		newIpv6, ipv6Err = requests.GetIP(cfg.Ipv6url)
	}
	if ipv4Err != nil {
		newIpv4 = ""
	}
	if ipv6Err != nil {
		newIpv6 = ""
	}

	return newIpv4, newIpv6, ipv4Err, ipv6Err
}
