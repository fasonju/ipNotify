package notify

import "log/slog"

func Notify(message string) {
	slog.Info(message)
}
