package download_test

import (
	"errors"
	"io"
	"scyllaDbAssignment/internal/cloud"
	"scyllaDbAssignment/internal/download"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDownloadVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := cloud.NewMockClientInterface(ctrl)

	fullVersion := "2025.1.1"
	versionSuffix := "2025.1"
	downloadableKey := "scylla-cqlsh-2025.1.1-0.20250407.1a896169dca9.x86_64.tar.gz"
	// when DownloadVersion("2025.1.0") then return "fake content"
	mockClient.EXPECT().
		DownloadVersion(versionSuffix, downloadableKey).
		Return(io.NopCloser(strings.NewReader("fake content")), nil)

	err := download.Run(mockClient, fullVersion, nil)

	assert.Nil(t, err)
}

func TestDownloadVersion_ClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := cloud.NewMockClientInterface(ctrl)
	returnedErr := errors.New("fake error")

	fullVersion := "2025.1.1"
	versionSuffix := "2025.1"
	downloadableKey := "scylla-cqlsh-2025.1.1-0.20250407.1a896169dca9.x86_64.tar.gz"
	// when DownloadVersion("2025.1.0") then return "fake content"
	mockClient.EXPECT().
		DownloadVersion(versionSuffix, downloadableKey).
		Return(nil, returnedErr)

	err := download.Run(mockClient, fullVersion, nil)

	assert.Error(t, err, "\"failed to download fullVersion fake error")
}
