package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"transitmanifest/app"
	"transitmanifest/domain"
	"transitmanifest/query"
)

func run(args []string, stdout, stderr io.Writer) error {
	config := app.FromEnvironment()
	application, err := app.Open(config)
	if err != nil {
		return err
	}
	defer application.Close()
	if len(args) == 0 {
		usage(stdout)
		return nil
	}
	switch args[0] {
	case "create":
		return create(application, args[1:], stdout)
	case "scan":
		return scanCommand(application, args[1:], stdout)
	case "close":
		return closeCommand(application, args[1:], stdout)
	case "handoff":
		return handoffCommand(application, args[1:], stdout)
	case "list":
		return writeJSON(stdout, application.Query.Manifests(queryFilter(args[1:])))
	case "report":
		return reportCommand(application, args[1:], stdout)
	case "audit":
		return auditCommand(application, args[1:], stdout)
	case "verify":
		return application.Audit.Verify()
	default:
		usage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func create(a *app.App, args []string, out io.Writer) error {
	set := flag.NewFlagSet("create", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	id := set.String("id", "", "manifest id")
	trip := set.String("trip", "", "trip id")
	plate := set.String("plate", "", "vehicle plate")
	origin := set.String("origin", "", "origin")
	destination := set.String("destination", "", "destination")
	if err := set.Parse(args); err != nil {
		return err
	}
	item, err := a.Manifests.Create(*id, *trip, *plate, *origin, *destination, a.Config.Actor, time.Now())
	if err != nil {
		return err
	}
	return writeJSON(out, item)
}

func scanCommand(a *app.App, args []string, out io.Writer) error {
	set := flag.NewFlagSet("scan", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	manifestID := set.String("manifest", "", "manifest id")
	scanID := set.String("scan", "", "scan id")
	packageID := set.String("package", "", "package id")
	station := set.String("station", "", "station")
	operation := set.String("operation", "", "load or unload")
	if err := set.Parse(args); err != nil {
		return err
	}
	accepted, exception, err := a.Scans.Submit(domain.PackageScan{ScanID: *scanID, ManifestID: *manifestID, PackageID: *packageID, Station: *station, Operator: a.Config.Actor, Operation: domain.ScanOperation(*operation), ScannedAt: time.Now()})
	if err != nil {
		return err
	}
	return writeJSON(out, map[string]any{"scan": accepted, "exception": exception})
}

func closeCommand(a *app.App, args []string, out io.Writer) error {
	set := flag.NewFlagSet("close", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	id := set.String("manifest", "", "manifest id")
	reason := set.String("reason", "", "reason")
	if err := set.Parse(args); err != nil {
		return err
	}
	item, err := a.Manifests.Close(*id, *reason, a.Config.Actor, time.Now())
	if err != nil {
		return err
	}
	return writeJSON(out, item)
}

func handoffCommand(a *app.App, args []string, out io.Writer) error {
	set := flag.NewFlagSet("handoff", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	id := set.String("manifest", "", "manifest id")
	from := set.String("from", "", "from station")
	to := set.String("to", "", "to station")
	signer := set.String("signer", "", "signer")
	note := set.String("note", "", "note")
	if err := set.Parse(args); err != nil {
		return err
	}
	receipt, err := a.Handoffs.Sign(*id, *from, *to, *signer, *note, time.Now())
	if err != nil {
		return err
	}
	return writeJSON(out, receipt)
}

func queryFilter(args []string) query.ManifestFilter {
	set := flag.NewFlagSet("list", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	status := set.String("status", "", "status")
	origin := set.String("origin", "", "origin")
	destination := set.String("destination", "", "destination")
	_ = set.Parse(args)
	return query.ManifestFilter{Status: domain.ManifestStatus(*status), Origin: *origin, Destination: *destination}
}

func reportCommand(a *app.App, args []string, out io.Writer) error {
	set := flag.NewFlagSet("report", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	day := set.String("date", time.Now().Format("2006-01-02"), "date")
	if err := set.Parse(args); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02", *day)
	if err != nil {
		return err
	}
	result, err := a.Reports.Daily(parsed)
	if err != nil {
		return err
	}
	return writeJSON(out, result)
}

func auditCommand(a *app.App, args []string, out io.Writer) error {
	set := flag.NewFlagSet("audit", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	id := set.String("manifest", "", "manifest id")
	if err := set.Parse(args); err != nil {
		return err
	}
	entries, err := a.Audit.List(*id, time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	return writeJSON(out, entries)
}
