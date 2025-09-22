package main

import (
	"flag"
	"fmt"
	"os"
	"scyllaDbAssignment/internal/cloud"
	"scyllaDbAssignment/internal/download"
	"scyllaDbAssignment/internal/listing"
	"scyllaDbAssignment/pkg"
	"scyllaDbAssignment/pkg/validation"
)

const errorMsgTemplate = "Error %s"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid syntax. Available commands are the following: list, download")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case pkg.ListCommand:
		runList()
	case pkg.DownloadCommand:
		runDownload()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

// runList handles flags and calls the listing package
func runList() {
	listFlags := flag.NewFlagSet("list", flag.ExitOnError)

	lt := listFlags.String(pkg.LesserThenFlag, "", "Filter items by substring")
	gt := listFlags.String(pkg.GreaterThenFlag, "", "Maximum number of items to list")

	if err := listFlags.Parse(os.Args[2:]); err != nil {
		fmt.Println("Failed to parse flags:", err)
		os.Exit(1)
	}

	ltValue := ""
	gtValue := ""
	if lt != nil && *lt != "" {
		ltValue = *lt
	}

	if gt != nil && *gt != "" {
		gtValue = *gt
	}

	err := validateListInput(gtValue, ltValue)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	listParams := listing.ListParams{
		LT: ltValue,
		GT: gtValue,
	}
	scyllaClient := cloud.New()
	versions, err := listing.Run(scyllaClient, listParams)
	if err != nil {
		fmt.Printf(errorMsgTemplate, err)
		os.Exit(2)
	}

	printVersions(versions)
}

func validateListInput(gt string, lt string) error {
	if lt != "" {
		return validation.ValidateVersion(lt)
	}

	if gt != "" {
		return validation.ValidateVersion(gt)
	}
	return nil
}

func runDownload() {
	downloadFlag := flag.NewFlagSet("download", flag.ExitOnError)
	output := downloadFlag.String(pkg.OutputFileFlag, "", "Path to save the downloaded file (default stdout)")
	if err := downloadFlag.Parse(os.Args[3:]); err != nil {
		fmt.Printf(errorMsgTemplate, err)
		os.Exit(2)
	}

	version := os.Args[2]
	err := validation.ValidateFullVersion(version)
	if err != nil {
		fmt.Printf(errorMsgTemplate, err)
		os.Exit(1)
	}

	scyllaClient := cloud.New()
	err = download.Run(scyllaClient, version, output)
	if err != nil {
		fmt.Printf(errorMsgTemplate, err)
		os.Exit(2)
	}
	fmt.Println("Download completed")
}

func printVersions(versions []pkg.Version) {
	fmt.Println("VERSION\tCLOUD")
	for _, v := range versions {
		fmt.Printf("%s\t%s\n", v.Name, v.CloudState)
	}
}
