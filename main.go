package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

// Compile-time version variable
var version string

// main is the entry point of the application. It initializes the configuration,
// parses command-line flags, and executes the certificate check process.
func main() {
	cfg := &Config{}
	setupFlags(cfg)
	flag.Parse()

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		flag.Usage()
		os.Exit(1)
	}

	if cfg.ShowVersion {
		fmt.Println("PEM Certificate Expiration Checker - Version:", version)
		return
	}

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// run coordinates the file discovery and processing based on the provided configuration.
// It determines whether to process a single file or a directory.
func run(cfg *Config) error {
	var files []string
	var err error

	if cfg.Dir != "" {
		files, err = getFilesFromDir(cfg.Dir, cfg.Ext)
		if err != nil {
			return fmt.Errorf("error reading directory: %w", err)
		}
	} else if cfg.File != "" {
		files = append(files, cfg.File)
	}

	reports := processFiles(files, cfg.Days)
	if err := displayReport(reports, cfg.ShowHeader); err != nil {
		return err
	}

	return nil
}

// processFiles handles the concurrent processing of multiple certificate files.
// It uses a WaitGroup and a channel to collect reports from goroutines.
func processFiles(files []string, days int) []CertificateReport {
	var wg sync.WaitGroup
	reportsChan := make(chan CertificateReport, len(files))

	for _, f := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			report, err := checkCertificate(file, days)
			if err != nil {
				reportsChan <- CertificateReport{FileName: file, Status: StatusInvalid, Error: err}
			} else {
				reportsChan <- report
			}
		}(f)
	}

	wg.Wait()
	close(reportsChan)

	var reports []CertificateReport
	for report := range reportsChan {
		reports = append(reports, report)
	}

	return reports
}
