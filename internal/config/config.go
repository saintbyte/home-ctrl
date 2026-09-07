package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Task represents a background task configuration
type Task struct {
	Name     string         `yaml:"name"`
	Schedule string         `yaml:"schedule"` // cron expression
	Enabled  bool           `yaml:"enabled"`
	Command  string         `yaml:"command"`
	Params   map[string]any `yaml:"params"`
}

// Widget represents a widget in the main view
type Widget struct {
	Name   string         `yaml:"name"`
	Params map[string]any `yaml:"params"`
}

// MainView represents the main view configuration
type MainView struct {
	Widgets []Widget `yaml:"widgets"`
}

// Config represents the application configuration
type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		TLS      bool   `yaml:"tls"`
		TLSCert  string `yaml:"tls_cert"`
		TLSKey   string `yaml:"tls_key"`
		CACert   string `yaml:"ca_cert"`
	} `yaml:"server"`

	Auth struct {
		Users      map[string]string `yaml:"users"`
		SessionTTL int               `yaml:"session_ttl_hours"`
	} `yaml:"auth"`

	DataDir string `yaml:"data_dir"`

	Tasks []Task `yaml:"tasks"`

	MainView MainView `yaml:"mainview"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			TLS      bool   `yaml:"tls"`
			TLSCert  string `yaml:"tls_cert"`
			TLSKey   string `yaml:"tls_key"`
			CACert   string `yaml:"ca_cert"`
		}{
			Host: "127.0.0.1",
			Port: 8080,
			TLS:  false,
		},
		Auth: struct {
			Users      map[string]string `yaml:"users"`
			SessionTTL int               `yaml:"session_ttl_hours"`
		}{
			Users: map[string]string{
				"admin": "admin123",
				"user":  "user123",
			},
			SessionTTL: 24, // 24 hours
		},
		Tasks: []Task{},
		MainView: MainView{
			Widgets: []Widget{},
		},
	}
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s %w", path, err)
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

// StringParam returns the string value of a task/widget parameter, or "" if absent
func StringParam(params map[string]any, name string) string {
	if v, ok := params[name].(string); ok {
		return v
	}
	return ""
}

// IntParam returns the int value of a task/widget parameter, or 0 if absent
func IntParam(params map[string]any, name string) int {
	switch v := params[name].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

// BoolParam returns the bool value of a task/widget parameter
func BoolParam(params map[string]any, name string) (bool, bool) {
	v, ok := params[name].(bool)
	return v, ok
}
