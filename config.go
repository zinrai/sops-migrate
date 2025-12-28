package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Target string       `yaml:"target"`
	Files  []FileConfig `yaml:"files"`
}

type FileConfig struct {
	Path      string `yaml:"path"`
	InputType string `yaml:"input_type"`
}

func LoadConfig(path string) (*Config, error) {
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

func (c *Config) GetInputType(path string) string {
	// pathからtargetプレフィックスを除去して比較
	relativePath := path
	if c.Target != "" {
		relativePath = strings.TrimPrefix(path, c.Target)
		relativePath = strings.TrimPrefix(relativePath, "/")
	}

	for _, f := range c.Files {
		if f.Path == relativePath {
			return f.InputType
		}
	}
	return ""
}
