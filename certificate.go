// Package main provides a utility to check the expiration status of PEM-encoded certificates.
package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Certificate statuses represent the possible states of a certificate.
const (
	// StatusValid indicates the certificate is not expired and not near expiration.
	StatusValid = "Valid"
	// StatusExpiring indicates the certificate is close to its expiration date.
	StatusExpiring = "Expiring"
	// StatusExpired indicates the certificate's expiration date has passed.
	StatusExpired = "Expired"
	// StatusInvalid indicates the file is not a valid PEM certificate.
	StatusInvalid = "Invalid"
)

// CertificateReport contains detailed information about a processed certificate,
// including its subject, issuer, and expiration status.
type CertificateReport struct {
	FileName     string    // Path to the certificate file
	CommonName   string    // Subject Common Name (CN)
	Issuer       string    // Issuer Common Name (CN)
	Status       string    // Current status (Valid, Expiring, Expired, or Invalid)
	Expiration   time.Time // Expiration date and time
	DaysToExpire int       // Number of days until expiration
	Error        error     // Any error encountered during processing
}

// checkCertificate reads a PEM-encoded certificate from the specified filePath,
// parses it, and determines its expiration status relative to expiringDays.
// It returns a CertificateReport containing the certificate's details.
func checkCertificate(filePath string, expiringDays int) (CertificateReport, error) {
	// Clean the path to prevent directory traversal (Gosec G304)
	cleanPath := filepath.Clean(filePath)

	// Read the raw file data
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return CertificateReport{FileName: cleanPath}, fmt.Errorf("could not read file: %w", err)
	}

	// Decode the PEM block and ensure it's a certificate
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return CertificateReport{FileName: cleanPath}, fmt.Errorf("not a valid PEM certificate")
	}

	// Parse the x509 certificate from the PEM block bytes
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return CertificateReport{FileName: cleanPath}, fmt.Errorf("could not parse certificate: %w", err)
	}

	// Calculate expiration details
	daysToExpire := int(time.Until(cert.NotAfter).Hours() / 24)
	status := StatusValid

	// Determine status based on expiration date
	if time.Now().After(cert.NotAfter) {
		status = StatusExpired
	} else if daysToExpire <= expiringDays {
		status = StatusExpiring
	}

	// Construct and return the report
	return CertificateReport{
		FileName:     cleanPath,
		CommonName:   cert.Subject.CommonName,
		Issuer:       cert.Issuer.CommonName,
		Status:       status,
		Expiration:   cert.NotAfter,
		DaysToExpire: daysToExpire,
	}, nil
}
