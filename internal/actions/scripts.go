package actions

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	SCRIPTS_PATH = "/scripts"
)

// / will put the arguments of the method as arugments to all scripts
func ExecuteScripts(previousIpv4 string, newIpv4 string, previousIpv6 string, newIpv6 string, message string) {
	files, err := os.ReadDir(SCRIPTS_PATH)
	if err != nil {
		slog.Error("Unable to read scripts folder", "error", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			scriptPath := filepath.Join(SCRIPTS_PATH, file.Name())

			slog.Info("Executing script", "name", file.Name())

			cmd := exec.Command(scriptPath, previousIpv4, newIpv4, previousIpv6, newIpv6, message)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				slog.Error("Failed to run script", "script", file.Name(), "error", err)
			}
		}
	}

	slog.Info("All scripts ran")
}
