package tui

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ConfigDirName = ".config/gitx"
	ConfigFileName = "config.toml"
	ConfigDirPath string
	ConfigFilePath string
	ConfigThemesDirPath string
)

func initializeConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	ConfigDirPath = filepath.Join(homeDir, ConfigDirName)
	ConfigFilePath = filepath.Join(ConfigDirPath, ConfigFileName)
	ConfigThemesDirPath = filepath.Join(ConfigDirPath, "themes")

	err = os.MkdirAll(ConfigDirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	err = os.MkdirAll(ConfigThemesDirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating themes directory: %w", err)
	}

	if _, err := os.Stat(ConfigFilePath); err != nil {
		if os.IsNotExist(err) {
			defaultConfig := fmt.Sprintf("Theme = %q\n", DefaultThemeName)
			if writeErr := os.WriteFile(ConfigFilePath, []byte(defaultConfig), 0644); writeErr != nil {
				return fmt.Errorf("failed to create default config file: %w", writeErr)
			}
		} else {
			return fmt.Errorf("failed to check config file: %w", err)
		}
	}

	return nil
}

func init() {
	if err := initializeConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
		os.Exit(1)
	}
}