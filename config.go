package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Files []FileConfig `yaml:"files"`
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
	for _, f := range c.Files {
		if f.Path == path {
			return f.InputType
		}
	}
	return ""
}
