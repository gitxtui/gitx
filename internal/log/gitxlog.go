package gitxlog

import (
	"log"
	"os"
	"path/filepath"
)

// SetupLogger sets up a logfile in the user's home directory.
// Also configures the standard logger to write the logs to the log file
func SetupLogger() (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	logFilePath := filepath.Join(home, ".gitx.log")

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.Println("--- Gitx session started ---")

	return file, nil
}
