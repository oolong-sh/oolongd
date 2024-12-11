package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func findConfigPath() (string, error) {
	var err error

	// read config file location from env var, fall back to ~/.config/oolong.toml
	configPath := os.Getenv("OOLONG_CONFIG_PATH")
	if configPath == "" {
		configPath, err = checkDefaultLocations()
		if err != nil {
			return "", err
		}
	}
	fmt.Println(configPath)

	configPath, err = expand(configPath)
	if err != nil {
		return "", err
	}

	return configPath, nil
}

func checkDefaultLocations() (string, error) {
	configPaths := []string{
		"~/.config/oolong/oolong.toml",
		"~/.config/oolong.toml",
		"./oolong.toml",
	}

	for _, p := range configPaths {
		configPath, err := expand(p)
		if err != nil {
			continue
		}

		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	return "", fmt.Errorf("No fallback config file found.")
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
