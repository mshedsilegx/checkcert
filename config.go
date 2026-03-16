package main

import (
	"flag"
	"fmt"
)

// Config holds the application's configuration
type Config struct {
	Dir         string
	File        string
	Ext         string
	Days        int
	ShowVersion bool
	ShowHeader  bool
}

// setupFlags configures the command-line flags for the application.
func setupFlags(cfg *Config) {
	flag.StringVar(&cfg.Dir, "dir", "", "Directory containing PEM files to scan")
	flag.StringVar(&cfg.File, "file", "", "Single PEM file to check")
	flag.StringVar(&cfg.Ext, "ext", ".pem", "Filter for PEM file extension (e.g., .crt, .pem)")
	flag.IntVar(&cfg.Days, "days", 30, "Number of days before expiration to flag as 'Expiring'")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Display the application version and exit")
	flag.BoolVar(&cfg.ShowHeader, "header", false, "Include a header in the output report")
}

// validateConfig ensures that the provided configuration is valid.
func validateConfig(cfg *Config) error {
	if !cfg.ShowVersion && cfg.Dir == "" && cfg.File == "" {
		return fmt.Errorf("either a directory or a file must be specified")
	}
	return nil
}
