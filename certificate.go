package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

// Certificate statuses
const (
	StatusValid    = "Valid"
	StatusExpiring = "Expiring"
	StatusExpired  = "Expired"
	StatusInvalid  = "Invalid"
)

// CertificateReport holds the details of the certificate
type CertificateReport struct {
	FileName     string
	CommonName   string
	Issuer       string
	Status       string
	Expiration   time.Time
	DaysToExpire int
	Error        error
}

func checkCertificate(filePath string, expiringDays int) (CertificateReport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return CertificateReport{FileName: filePath}, fmt.Errorf("could not read file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return CertificateReport{FileName: filePath}, fmt.Errorf("not a valid PEM certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return CertificateReport{FileName: filePath}, fmt.Errorf("could not parse certificate: %w", err)
	}

	daysToExpire := int(time.Until(cert.NotAfter).Hours() / 24)
	status := StatusValid
	if time.Now().After(cert.NotAfter) {
		status = StatusExpired
	} else if daysToExpire <= expiringDays {
		status = StatusExpiring
	}

	return CertificateReport{
		FileName:     filePath,
		CommonName:   cert.Subject.CommonName,
		Issuer:       cert.Issuer.CommonName,
		Status:       status,
		Expiration:   cert.NotAfter,
		DaysToExpire: daysToExpire,
	}, nil
}
