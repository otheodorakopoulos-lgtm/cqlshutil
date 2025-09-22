package cloud

import (
	"io"
	"scyllaDbAssignment/pkg"
)

type ClientInterface interface {
	ListVersions() ([]pkg.Version, error)
	DownloadVersion(versionSuffix, downloadableKey string) (io.ReadCloser, error)
}
