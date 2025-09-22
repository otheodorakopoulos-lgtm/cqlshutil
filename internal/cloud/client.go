package cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scyllaDbAssignment/pkg"
	"time"
)

//go:generate mockgen -source=clientinterface.go -destination=mockclient.go -package=cloud

type scyllaVersionItem struct {
	ID          int    `json:"id"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NewCluster  string `json:"newCluster"` // "ENABLED", "DISABLED"
}

type scyllaVersionData struct {
	ScyllaVersions []scyllaVersionItem `json:"scyllaVersions"`
}

type scyllaVersionResponse struct {
	Data scyllaVersionData `json:"data"`
}

const (
	listScyllaVersionsPath     = "https://api.cloud.scylladb.com/deployment/scylla-versions"
	downloadScyllaVersionPath  = "https://downloads.scylladb.com/downloads/scylla/relocatable/scylladb-%s/%s"
	listVersionsTimeoutSeconds = 10
	downloadTimeoutSeconds     = 30
)

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{http: &http.Client{Timeout: 60 * time.Second}}
}

func (c Client) ListVersions() ([]pkg.Version, error) {
	client := &http.Client{Timeout: listVersionsTimeoutSeconds * time.Second}

	req, err := http.NewRequest("GET", listScyllaVersionsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad response: %s - %s", resp.Status, string(body))
	}

	var apiResp scyllaVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response JSON failed: %w", err)
	}

	versions := make([]pkg.Version, 0, len(apiResp.Data.ScyllaVersions))
	for _, item := range apiResp.Data.ScyllaVersions {
		cloudState := item.NewCluster
		if cloudState == "" {
			cloudState = "N/A"
		}
		versions = append(versions, pkg.Version{
			Name:       item.Version,
			CloudState: cloudState,
		})
	}

	return versions, nil
}

func (c Client) DownloadVersion(versionSuffix, downloadableKey string) (io.ReadCloser, error) {
	url := fmt.Sprintf(
		downloadScyllaVersionPath,
		versionSuffix,
		downloadableKey)

	client := &http.Client{Timeout: downloadTimeoutSeconds * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("bad status downloading %s: %s", url, resp.Status)
	}

	return resp.Body, nil
}
