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

// ExecuteScripts runs all executable files in the scripts directory,
// passing IP change details and a message as arguments.
func ExecuteScripts(previousIpv4 string, newIpv4 string, previousIpv6 string, newIpv6 string, message string) {
	files, err := os.ReadDir(SCRIPTS_PATH)
	if err != nil {
		slog.Error("Unable to read scripts folder", "error", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			slog.Warn("Unable to check file permissions, continuing without checking")
		} else {
			// Ensure the script is executable by the user
			if fileInfo.Mode()&0111 == 0 {
				newMode := fileInfo.Mode() | 0100
				scriptPath := filepath.Join(SCRIPTS_PATH, file.Name())
				if err := os.Chmod(scriptPath, newMode); err != nil {
					slog.Error("Failed to make script executable", "script", file.Name(), "error", err)
					continue
				}
				slog.Info("Made script executable", "script", file.Name())
			}
		}

		// Run the script with the IPs and message as arguments
		scriptPath := filepath.Join(SCRIPTS_PATH, file.Name())
		slog.Info("Executing script", "name", file.Name())

		cmd := exec.Command(scriptPath, previousIpv4, newIpv4, previousIpv6, newIpv6, message)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			slog.Error("Failed to run script", "script", file.Name(), "error", err)
		}
	}

	slog.Info("All scripts ran")
}
