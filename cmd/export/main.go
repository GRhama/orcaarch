package main

import (
	"fmt"
	"os"

	"orcaarch/internal/export"
	"orcaarch/internal/report"
	"orcaarch/internal/service"
)

func main() {
	result, err := service.Process(1000)
	if err != nil {
		panic(err)
	}
	write("reconciliation.csv", func(f *os.File) error {
		return export.WriteReconciliationCSV(f, report.ToReconciliationRows(result.Reconciliations))
	})
	write("ledger.csv", func(f *os.File) error {
		return export.WriteLedgerCSV(f, report.ToLedgerRows(result.LedgerEntries))
	})
	write("risk.csv", func(f *os.File) error {
		return export.WriteRiskCSV(f, report.ToRiskRows(result.Quarantine))
	})
}

func write(name string, fn func(*os.File) error) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	if err := fn(f); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", name, err)
		os.Exit(1)
	}
	fmt.Println("wrote", name)
}
