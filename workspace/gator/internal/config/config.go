package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

// SetUser sets the current_user_name field and writes the config to disk.
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUser = userName
	return write(*cfg)
}

// Read loads and parses the configuration file from the user's home directory.
func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("Error finding path: %w", err)
	}

	// Create path
	configPath := filepath.Join(home, ".gatorconfig.json")

	// Load raw bytes
	data, err := os.ReadFile(configPath)

	// Unmarshall the data
	var configs Config
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding data: %w", err)
	}

	return configs, nil
}

// Helper function: write
func write(cfg Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	configPath := filepath.Join(home, ".gatorconfig.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
