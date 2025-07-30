package notify

import (
	"fmt"
	"log/slog"
	"net/smtp"

	"github.com/fasonju/ipNotify/internal/types"
)

func NotifySMTP(cfg *types.Config, message string) {
	subject := "Subject: Ip update notification\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"

	hostAddr := cfg.SmtpServer + ":" + fmt.Sprint(cfg.SmtpPort)
	slog.Info("Sending Email", "mail", message)
	auth := smtp.PlainAuth("", cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpServer)

	mail := []byte(subject + mime + "\r\n" + message)

	err := smtp.SendMail(hostAddr, auth, cfg.SmtpFrom, []string{cfg.SmtpTo}, mail)
	if err != nil {
		slog.Error("Error during sending of message", "error", err)
	} else {
		slog.Info("Email sent", "recipient", cfg.SmtpTo)
	}
}
