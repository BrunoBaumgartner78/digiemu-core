package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func runExport(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "export subcommands: unit")
		os.Exit(2)
	}

	switch args[0] {
	case "unit":
		fs := flag.NewFlagSet("export unit", flag.ExitOnError)
		unitKey := fs.String("unit", "", "unit key (required)")
		data := fs.String("data", "./data", "data directory")
		withAudit := fs.Bool("audit", false, "include audit events for this unit")
		pretty := fs.Bool("pretty", false, "pretty-print JSON")
		fs.Parse(args[1:])

		if *unitKey == "" {
			fmt.Fprintln(os.Stderr, "--unit is required")
			fs.Usage()
			os.Exit(2)
		}

		repo := fsrepo.NewUnitRepo(*data)

		var audit ports.AuditLogByUnitReader
		if *withAudit {
			audit = fsrepo.NewAuditByUnitReader(*data)
		}

		uc := usecases.ExportUnitSnapshot{Repo: repo, Audit: audit}

		out, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{
			UnitKey:      *unitKey,
			IncludeAudit: *withAudit,
		})
		if err != nil {
			log.Fatalf("export unit: %v", err)
		}

		var b []byte
		if *pretty {
			b, err = json.MarshalIndent(out, "", "  ")
		} else {
			b, err = json.Marshal(out)
		}
		if err != nil {
			log.Fatalf("export marshal: %v", err)
		}

		fmt.Println(string(b))

	default:
		fmt.Fprintln(os.Stderr, "export subcommands: unit")
		os.Exit(2)
	}
}
