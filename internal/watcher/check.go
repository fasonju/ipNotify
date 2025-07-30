package watcher

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/fasonju/ipNotify/internal/actions"
	"github.com/fasonju/ipNotify/internal/types"
)

// checkIpDiffAndNotify compares current IPs with previous ones,
// triggers notifications if changes are detected.
func checkIpDiffAndNotify(previousIpv4, previousIpv6 string, cfg *types.Config) (string, string, error) {
	newIpv4, newIpv6, ipv4Err, ipv6Err := fetchCurrentIPs(cfg)
	if ipv4Err != nil && ipv6Err != nil {
		return previousIpv4, previousIpv6, ipv4Err
	}

	if ipsChanged(previousIpv4, newIpv4, previousIpv6, newIpv6) {
		message := formatChangeMessage(previousIpv4, newIpv4, previousIpv6, newIpv6, cfg, ipv4Err, ipv6Err)

		if cfg.SmtpEnabled {
			slog.Info("Notifying through SMTP")
			actions.NotifySMTP(cfg, message)
		}

		if cfg.ScriptsEnabled {
			slog.Info("Executing scripts")
			actions.ExecuteScripts(previousIpv4, newIpv4, previousIpv6, newIpv6, message)
		}
	} else {
		slog.Info("No IP change", "ipv4", newIpv4, "ipv6", newIpv6)
	}

	if ipv4Err != nil {
		return newIpv4, newIpv6, ipv4Err
	}
	if ipv6Err != nil {
		return newIpv4, newIpv6, ipv4Err
	}

	return newIpv4, newIpv6, nil
}

// ipsChanged returns true if either IPv4 or IPv6 has changed.
func ipsChanged(prev4, new4, prev6, new6 string) bool {
	return new4 != prev4 || new6 != prev6
}

// formatChangeMessage builds a message describing the IP changes.
func formatChangeMessage(prev4, new4, prev6, new6 string, cfg *types.Config, ipv4Err, ipv6Err error) string {
	var builder strings.Builder

	if ipv4Err != nil {
		builder.WriteString(fmt.Sprintf("IPV4 IP query error: %s", ipv4Err.Error()))
	}
	if ipv6Err != nil {
		builder.WriteString(fmt.Sprintf("IPV6 IP query error: %s", ipv6Err.Error()))
	}

	if cfg.Ipv4Enabled && new4 != prev4 {
		builder.WriteString(fmt.Sprintf("IPv4 changed from %s to %s\n", prev4, new4))
	}
	if cfg.Ipv6Enabled && new6 != prev6 {
		builder.WriteString(fmt.Sprintf("IPv6 changed from %s to %s\n", prev6, new6))
	}

	return builder.String()
}
