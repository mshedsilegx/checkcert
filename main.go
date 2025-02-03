package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Compile-time version variable
var version string

// CertificateReport holds the details of the certificate
type CertificateReport struct {
	FileName     string
	CommonName   string
	Issuer       string
	Status       string
	Expiration   time.Time
	DaysToExpire int
}

func main() {
	// Command-line flags
	dir := flag.String("dir", "", "Directory containing PEM files")
	file := flag.String("file", "", "Single PEM file")
	ext := flag.String("ext", ".pem", "PEM file extension to search for")
	days := flag.Int("days", 30, "Number of days to check for expiry")
	showVersion := flag.Bool("version", false, "Display application version")
	showHeader := flag.Bool("header", false, "Display report header")

	flag.Parse()

	if *showVersion {
		fmt.Println("PEM Certificate Expiration Checker - Version:", version)
		return
	}

	var files []string
	var err error

	if *dir != "" {
		files, err = getFilesFromDir(*dir, *ext)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}
	} else if *file != "" {
		files = append(files, *file)
	} else {
		fmt.Println("Either a directory or a file must be specified.")
		return
	}

	reports := []CertificateReport{}
	for _, f := range files {
		report := checkCertificate(f, *days)
		reports = append(reports, report)
	}

	displayReport(reports, *showHeader)
}

func getFilesFromDir(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func checkCertificate(filePath string, expiringDays int) CertificateReport {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return CertificateReport{FileName: filePath, Status: "Invalid"}
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return CertificateReport{FileName: filePath, Status: "Invalid"}
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return CertificateReport{FileName: filePath, Status: "Invalid"}
	}

	daysToExpire := int(time.Until(cert.NotAfter).Hours() / 24)
	status := "Valid"
	if time.Now().After(cert.NotAfter) {
		status = "Expired"
	} else if daysToExpire <= expiringDays {
		status = "Expiring"
	}

	return CertificateReport{
		FileName:     filePath,
		CommonName:   cert.Subject.CommonName,
		Issuer:       cert.Issuer.CommonName,
		Status:       status,
		Expiration:   cert.NotAfter,
		DaysToExpire: daysToExpire,
	}
}

func displayReport(reports []CertificateReport, displayHeader bool) {
	if displayHeader {
		fmt.Printf("%-60s %-40s %-60s %-10s %-15s %-5s\n", "Object Name (PEM File)", "Common Name (CN)", "Issuer", "Status", "Expiration", "Days")
		fmt.Println(strings.Repeat("-", 200))
	}
	for _, report := range reports {
		fmt.Printf("%-60s %-40s %-60s %-10s %-15s %-5d\n",
			report.FileName, report.CommonName, report.Issuer, report.Status, report.Expiration.Format("2006-01-02"), report.DaysToExpire)
	}
}
