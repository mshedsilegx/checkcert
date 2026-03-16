package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// displayReport prints a formatted table of certificate reports to standard output.
// It uses a tabwriter to ensure columns are properly aligned.
func displayReport(reports []CertificateReport, displayHeader bool) error {
	// Initialize a new tabwriter to format the output into columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Print the header if requested
	if displayHeader {
		if _, err := fmt.Fprintln(w, "Object Name (PEM File)\tCommon Name (CN)\tIssuer\tStatus\tExpiration\tDays\tError"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "----------------------\t----------------\t------\t------\t----------\t----\t-----"); err != nil {
			return err
		}
	}

	// Iterate through each report and print its details
	for _, report := range reports {
		// Handle potential error messages
		errMsg := ""
		if report.Error != nil {
			errMsg = report.Error.Error()
		}

		// Format the expiration date if it exists
		expirationDate := ""
		if !report.Expiration.IsZero() {
			expirationDate = report.Expiration.Format("2006-01-02")
		}

		// Print the formatted row
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			report.FileName,
			report.CommonName,
			report.Issuer,
			report.Status,
			expirationDate,
			report.DaysToExpire,
			errMsg,
		); err != nil {
			return err
		}
	}

	// Flush the tabwriter to output the buffered data
	return w.Flush()
}
