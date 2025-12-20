package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	WelcomeMessage string `json:"welcome_message"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".freeport_config.json"
	}

	return filepath.Join(home, ".freeport_config.json")
}

func Load() *Config {
	cfg := &Config{
		WelcomeMessage: "Welcome to Freeport!",
	}

	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return cfg
	}

	json.Unmarshal(data, cfg)
	return cfg
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(getConfigPath(), data, 0644)
}