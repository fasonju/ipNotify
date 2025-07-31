package watcher

import (
	"log/slog"
	"sync"
	"time"

	"github.com/fasonju/ipNotify/internal/requests"
	"github.com/fasonju/ipNotify/internal/types"
)

const (
	maxRetries = 3
	retryDelay = 500 * time.Millisecond
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

// retryGetIP tries to get the IP from the given URL up to maxRetries times with retryDelay between attempts.
func retryGetIP(url string, maxRetries int, retryDelay time.Duration) (string, error) {
	var ip string
	var err error
	for range maxRetries {
		ip, err = requests.GetIP(url)
		if err == nil {
			return ip, nil
		}
		time.Sleep(retryDelay)
	}
	return "", err
}

// fetchCurrentIPs returns the current IPv4 and IPv6 addresses based on the config.
// It returns separate errors for IPv4 and IPv6 lookups if they fail.
// If an IP version is disabled, its address will be empty and error nil.
func fetchCurrentIPs(cfg *types.Config) (string, string, error, error) {
	var newIpv4, newIpv6 string
	var ipv4Err, ipv6Err error
	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	// IPv4 goroutine
	go func() {
		defer waitGroup.Done()
		if cfg.Ipv4Enabled {
			newIpv4, ipv4Err = retryGetIP(cfg.Ipv4url, maxRetries, retryDelay)
			if ipv4Err != nil {
				newIpv4 = ""
			}
		}
	}()

	// IPv6 goroutine
	go func() {
		defer waitGroup.Done()
		if cfg.Ipv6Enabled {
			newIpv6, ipv6Err = retryGetIP(cfg.Ipv6url, maxRetries, retryDelay)
			if ipv6Err != nil {
				newIpv6 = ""
			}
		}
	}()

	waitGroup.Wait()
	return newIpv4, newIpv6, ipv4Err, ipv6Err
}
