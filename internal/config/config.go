package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	
	Auth struct {
		Users     map[string]string `yaml:"users"`
		SessionTTL int              `yaml:"session_ttl_hours"`
	} `yaml:"auth"`
	
	DataDir string `yaml:"data_dir"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		}{
			Host: "127.0.0.1",
			Port: 8080,
		},
		Auth: struct {
			Users     map[string]string `yaml:"users"`
			SessionTTL int              `yaml:"session_ttl_hours"`
		}{
			Users: map[string]string{
				"admin": "admin123",
				"user":  "user123",
			},
			SessionTTL: 24, // 24 hours
		},
		DataDir: "data",
	}
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {
	// If config file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetServerAddress returns the server address (host:port)
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}