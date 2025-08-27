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

func setupFlags(cfg *Config) {
	flag.StringVar(&cfg.Dir, "dir", "", "Directory containing PEM files")
	flag.StringVar(&cfg.File, "file", "", "Single PEM file")
	flag.StringVar(&cfg.Ext, "ext", ".pem", "PEM file extension to search for")
	flag.IntVar(&cfg.Days, "days", 30, "Number of days to check for expiry")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Display application version")
	flag.BoolVar(&cfg.ShowHeader, "header", false, "Display report header")
}

func validateConfig(cfg *Config) error {
	if !cfg.ShowVersion && cfg.Dir == "" && cfg.File == "" {
		return fmt.Errorf("either a directory or a file must be specified")
	}
	return nil
}
