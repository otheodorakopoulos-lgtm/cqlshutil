package download

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"scyllaDbAssignment/internal/cloud"

	"scyllaDbAssignment/pkg"
	"strings"
)

var downloadables = []pkg.Downloadable{
	// 2025.1.x releases
	{"scylla-cqlsh-2025.1.0-0.20250325.9dca28d2b818.aarch64.tar.gz", "2025.1.0"},
	{"scylla-cqlsh-2025.1.0-0.20250325.9dca28d2b818.x86_64.tar.gz", "2025.1.0"},
	{"scylla-cqlsh-2025.1.0~rc0-0.20250128.f407799f252d.aarch64.tar.gz", "2025.1.0~rc0"},
	{"scylla-cqlsh-2025.1.0~rc0-0.20250128.f407799f252d.x86_64.tar.gz", "2025.1.0~rc0"},
	{"scylla-cqlsh-2025.1.0~rc1-0.20250202.28b889668011.aarch64.tar.gz", "2025.1.0~rc1"},
	{"scylla-cqlsh-2025.1.0~rc1-0.20250202.28b889668011.x86_64.tar.gz", "2025.1.0~rc1"},
	{"scylla-cqlsh-2025.1.0~rc2-0.20250216.6ee17795783f.aarch64.tar.gz", "2025.1.0~rc2"},
	{"scylla-cqlsh-2025.1.0~rc2-0.20250216.6ee17795783f.x86_64.tar.gz", "2025.1.0~rc2"},
	{"scylla-cqlsh-2025.1.0~rc3-0.20250223.aa5cb15166d3.aarch64.tar.gz", "2025.1.0~rc3"},
	{"scylla-cqlsh-2025.1.0~rc3-0.20250223.aa5cb15166d3.x86_64.tar.gz", "2025.1.0~rc3"},
	{"scylla-cqlsh-2025.1.0~rc4-0.20250323.bc983017832c.aarch64.tar.gz", "2025.1.0~rc4"},
	{"scylla-cqlsh-2025.1.0~rc4-0.20250323.bc983017832c.x86_64.tar.gz", "2025.1.0~rc4"},
	{"scylla-cqlsh-2025.1.1-0.20250407.1a896169dca9.aarch64.tar.gz", "2025.1.1"},
	{"scylla-cqlsh-2025.1.1-0.20250407.1a896169dca9.x86_64.tar.gz", "2025.1.1"},
	{"scylla-cqlsh-2025.1.2-0.20250422.502c62d91d48.aarch64.tar.gz", "2025.1.2"},
	{"scylla-cqlsh-2025.1.2-0.20250422.502c62d91d48.x86_64.tar.gz", "2025.1.2"},
	{"scylla-cqlsh-2025.1.3-0.20250529.5a67119dce9d.aarch64.tar.gz", "2025.1.3"},
	{"scylla-cqlsh-2025.1.3-0.20250529.5a67119dce9d.x86_64.tar.gz", "2025.1.3"},
	{"scylla-cqlsh-2025.1.4-0.20250707.20afd2776561.aarch64.tar.gz", "2025.1.4"},
	{"scylla-cqlsh-2025.1.4-0.20250707.20afd2776561.x86_64.tar.gz", "2025.1.4"},
	{"scylla-cqlsh-2025.1.5-0.20250715.23bcff5c5279.aarch64.tar.gz", "2025.1.5"},
	{"scylla-cqlsh-2025.1.5-0.20250715.23bcff5c5279.x86_64.tar.gz", "2025.1.5"},
	{"scylla-cqlsh-2025.1.6-0.20250811.38f4c3325d1f.aarch64.tar.gz", "2025.1.6"},
	{"scylla-cqlsh-2025.1.6-0.20250811.38f4c3325d1f.x86_64.tar.gz", "2025.1.6"},
	{"scylla-cqlsh-2025.1.7-0.20250909.32471fa8db82.aarch64.tar.gz", "2025.1.7"},
	{"scylla-cqlsh-2025.1.7-0.20250909.32471fa8db82.x86_64.tar.gz", "2025.1.7"},

	// 2025.2.x releases
	{"scylla-cqlsh-2025.2.0-0.20250625.33e947e75342.aarch64.tar.gz", "2025.2.0"},
	{"scylla-cqlsh-2025.2.0-0.20250625.33e947e75342.x86_64.tar.gz", "2025.2.0"},
	{"scylla-cqlsh-2025.2.0~rc0-0.20250507.b3dbfaf27a47.aarch64.tar.gz", "2025.2.0~rc0"},
	{"scylla-cqlsh-2025.2.0~rc0-0.20250507.b3dbfaf27a47.x86_64.tar.gz", "2025.2.0~rc0"},
	{"scylla-cqlsh-2025.2.0~rc1-0.20250513.6f1efcff315f.aarch64.tar.gz", "2025.2.0~rc1"},
	{"scylla-cqlsh-2025.2.0~rc1-0.20250513.6f1efcff315f.x86_64.tar.gz", "2025.2.0~rc1"},
	{"scylla-cqlsh-2025.2.1-0.20250716.7bb43d812e1c.aarch64.tar.gz", "2025.2.1"},
	{"scylla-cqlsh-2025.2.1-0.20250716.7bb43d812e1c.x86_64.tar.gz", "2025.2.1"},
	{"scylla-cqlsh-2025.2.2-0.20250808.d845de84aa84.aarch64.tar.gz", "2025.2.2"},
	{"scylla-cqlsh-2025.2.2-0.20250808.d845de84aa84.x86_64.tar.gz", "2025.2.2"},
}

// Run downloads the specified version of ScyllaDB relocatable package.
// Example: Run("2025.1.2")
func Run(scyllaClient cloud.ClientInterface, fullVersion string, outputFile *string) error {
	if scyllaClient == nil {
		return errors.New("scylla client is nil")
	}
	parts := strings.Split(fullVersion, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid fullVersion: %s", fullVersion)
	}
	versionSuffix := parts[0] + "." + parts[1]

	downloadableKey, err := selectDownloadableKey(fullVersion)
	if err != nil {
		return fmt.Errorf("could not find downloadable file: %w", err)
	}

	body, err := scyllaClient.DownloadVersion(versionSuffix, *downloadableKey)
	if err != nil {
		return fmt.Errorf("failed to download fullVersion %w", err)
	}

	return writeToOutput(body, outputFile)
}

func writeToOutput(body io.ReadCloser, outputFile *string) error {
	var out io.Writer
	if outputFile != nil && *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	if _, err := io.Copy(out, body); err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}
	return nil
}

func selectDownloadableKey(version string) (*string, error) {
	archSuffix := ""
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		archSuffix = "x86_64"
	case "arm64":
		archSuffix = "aarch64"
	default:
		archSuffix = "x86_64"
	}

	for _, item := range downloadables {
		if item.Version == version && strings.Contains(item.FullName, archSuffix) {
			return &item.FullName, nil
		}
	}
	return nil, fmt.Errorf("no downloadable found for version %s and architecture %s", version, arch)
}
