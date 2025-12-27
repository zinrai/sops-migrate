package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	dryRun := flag.Bool("dry-run", false, "Show commands without executing")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -config is required")
		flag.Usage()
		os.Exit(1)
	}

	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	migrator := NewMigrator(config, *dryRun)
	results, err := migrator.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printSummary(results, *dryRun)

	if hasFailure(results) {
		os.Exit(1)
	}
}

func printSummary(results []Result, dryRun bool) {
	if dryRun {
		fmt.Printf("\nDry run: %d files\n", len(results))
		return
	}

	succeeded := 0
	var failedFiles []string

	for _, r := range results {
		if r.Success {
			succeeded++
		} else {
			failedFiles = append(failedFiles, r.Path)
		}
	}

	failed := len(failedFiles)
	fmt.Printf("\nCompleted: %d succeeded, %d failed\n", succeeded, failed)

	if failed > 0 {
		fmt.Println("\nFailed files:")
		for _, f := range failedFiles {
			fmt.Printf("  %s\n", f)
		}
	}
}

func hasFailure(results []Result) bool {
	for _, r := range results {
		if !r.Success {
			return true
		}
	}
	return false
}
