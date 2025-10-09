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

	_, err = os.Stat(ConfigFilePath)
	if err == nil {
		// File exists
		// fmt.Println("Config file exists:", ConfigFilePath)
	} else if os.IsNotExist(err) {
    	// File does not exist
		defaultConfigContent := fmt.Sprintf("Theme = \"%s\"\n", DefaultThemeName)
    	err = os.WriteFile(ConfigFilePath, []byte(defaultConfigContent), 0644)
		if err != nil {
			return fmt.Errorf("error creating default config file: %w", err)
		}
	} else {
		// Some other error
		return fmt.Errorf("error checking config file: %w", err)
	}

	return nil
}

func init() {
	if err := initializeConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
		os.Exit(1)
	}
}