package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func displayReport(reports []CertificateReport, displayHeader bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	if displayHeader {
		fmt.Fprintln(w, "Object Name (PEM File)\tCommon Name (CN)\tIssuer\tStatus\tExpiration\tDays\tError")
		fmt.Fprintln(w, "----------------------\t----------------\t------\t------\t----------\t----\t-----")
	}

	for _, report := range reports {
		errMsg := ""
		if report.Error != nil {
			errMsg = report.Error.Error()
		}

		expirationDate := ""
		if !report.Expiration.IsZero() {
			expirationDate = report.Expiration.Format("2006-01-02")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			report.FileName,
			report.CommonName,
			report.Issuer,
			report.Status,
			expirationDate,
			report.DaysToExpire,
			errMsg,
		)
	}
}
