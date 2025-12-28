package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Migrator struct {
	config *Config
	dryRun bool
}

type Result struct {
	Path    string
	Command string
	Success bool
}

func NewMigrator(config *Config, dryRun bool) *Migrator {
	return &Migrator{
		config: config,
		dryRun: dryRun,
	}
}

func (m *Migrator) Run() ([]Result, error) {
	if !m.dryRun {
		if _, err := exec.LookPath("sops"); err != nil {
			return nil, fmt.Errorf("sops command not found: %w", err)
		}
	}

	files, err := m.collectFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to collect files: %w", err)
	}

	var results []Result
	for _, file := range files {
		result := m.processFile(file)
		results = append(results, result)
	}

	return results, nil
}

func (m *Migrator) collectFiles() ([]string, error) {
	var files []string

	target := m.config.Target
	if target == "" {
		target = "."
	}

	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (m *Migrator) processFile(path string) Result {
	args := m.buildArgs(path)
	command := "sops " + strings.Join(args, " ")

	result := Result{
		Path:    path,
		Command: command,
		Success: true,
	}

	if m.dryRun {
		fmt.Println(command)
		return result
	}

	cmd := exec.Command("sops", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Success = false
		fmt.Printf("failed: %s\n", command)
	} else {
		fmt.Printf("ok: %s\n", command)
	}

	return result
}

func (m *Migrator) buildArgs(path string) []string {
	args := []string{"encrypt"}

	inputType := m.config.GetInputType(path)
	if inputType != "" {
		args = append(args, "--input-type", inputType)
	}

	args = append(args, "-i", path)

	return args
}
