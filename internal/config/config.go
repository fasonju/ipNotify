package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fasonju/ipNotify/internal/types"
	"github.com/joho/godotenv"
)

const (
	IPV4_URL = "https://ipv4.icanhazip.com"
	IPV6_URL = "https://ipv6.icanhazip.com"
)

func LoadConfig() (*types.Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, fmt.Errorf("failed to load vars from .env file: %w", err)
		}
		slog.Info(".env file loaded")
	}

	ipv4Enabled := os.Getenv("IPV4_ENABLED") == "true"
	ipv6Enabled := os.Getenv("IPV6_ENABLED") == "true"
	if !ipv4Enabled && !ipv6Enabled {
		return nil, fmt.Errorf("both IPV4 and IPV6 are disabled")
	}

	intervalStr := os.Getenv("INTERVAL")
	if intervalStr == "" {
		return nil, fmt.Errorf("INTERVAL is not set")
	}
	interval, err := parseInterval(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid INTERVAL: %w", err)
	}
	smtpEnabled := os.Getenv("SMTP_ENABLED") == "true"
	var (
		smtpServer   = os.Getenv("SMTP_SERVER")
		smtpUsername = os.Getenv("SMTP_USERNAME")
		smtpPassword = os.Getenv("SMTP_PASSWORD")
		smtpFrom     = os.Getenv("SMTP_FROM")
		smtpTo       = os.Getenv("SMTP_TO")
		smtpPortStr  = os.Getenv("SMTP_PORT")
		smtpPort     int
	)

	if smtpEnabled {
		missing := []string{}
		if smtpServer == "" {
			missing = append(missing, "SMTP_SERVER")
		}
		if smtpUsername == "" {
			missing = append(missing, "SMTP_USERNAME")
		}
		if smtpPassword == "" {
			missing = append(missing, "SMTP_PASSWORD")
		}
		if smtpFrom == "" {
			missing = append(missing, "SMTP_FROM")
		}
		if smtpTo == "" {
			missing = append(missing, "SMTP_TO")
		}
		if smtpPortStr == "" {
			missing = append(missing, "SMTP_PORT")
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("missing required SMTP environment variables: %s", strings.Join(missing, ", "))
		}

		smtpPort, err = strconv.Atoi(smtpPortStr)
		if err != nil || smtpPort <= 0 {
			return nil, fmt.Errorf("invalid SMTP_PORT: %s", smtpPortStr)
		}
	}

	return &types.Config{
		Ipv4url:      IPV4_URL,
		Ipv6url:      IPV6_URL,
		Ipv4Enabled:  ipv4Enabled,
		Ipv6Enabled:  ipv6Enabled,
		Interval:     interval,
		SmtpEnabled:  smtpEnabled,
		SmtpServer:   smtpServer,
		SmtpUsername: smtpUsername,
		SmtpPassword: smtpPassword,
		SmtpFrom:     smtpFrom,
		SmtpTo:       smtpTo,
		SmtpPort:     smtpPort,
	}, nil
}

func parseInterval(interval string) (time.Duration, error) {
	interval = strings.TrimSpace(interval)
	if len(interval) < 2 {
		return 0, fmt.Errorf("not a valid interval: %s", interval)
	}
	unit := interval[len(interval)-1]
	numPart := interval[:len(interval)-1]

	value, err := strconv.Atoi(numPart)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 's':
		return time.Duration(value) * time.Second, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	case 'd':
		return time.Duration(value) * 24 * time.Hour, nil
	case 'w':
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case 'y':
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown time unit: %c", unit)
	}
}
