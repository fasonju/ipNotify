package types

import "time"

type Config struct {
	Ipv4url        string
	Ipv6url        string
	Interval       time.Duration
	Ipv4Enabled    bool
	Ipv6Enabled    bool
	SmtpEnabled    bool
	SmtpServer     string
	SmtpUsername   string
	SmtpPassword   string
	SmtpFrom       string
	SmtpTo         string
	SmtpPort       int
	ScriptsEnabled bool
}
