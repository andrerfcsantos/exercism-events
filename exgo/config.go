package exgo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type UserConfig struct {
	Apibaseurl string `json:"apibaseurl"`
	Token      string `json:"token"`
	Workspace  string `json:"workspace"`
}

func cliConfigDir() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir != "" {
			return filepath.Join(dir, "exercism")
		}
	} else {
		dir := os.Getenv("EXERCISM_CONFIG_HOME")
		if dir != "" {
			return dir
		}
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = filepath.Join(os.Getenv("HOME"), ".config")
		}
		if dir != "" {
			return filepath.Join(dir, "exercism")
		}
	}
	// If all else fails, use the current directory.
	dir, _ = os.Getwd()
	return dir
}

func findToken() (string, error) {
	token := os.Getenv("EXERCISM_TOKEN")

	if token != "" {
		return token, nil
	}

	configDir := cliConfigDir()
	configFile := filepath.Join(configDir, "user.json")

	// open config file and parse json
	f, err := os.Open(configFile)
	if err != nil {
		return "", fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	var config UserConfig
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return "", fmt.Errorf("decoding config file: %w", err)
	}

	return config.Token, nil
}
