package listing

import (
	"errors"
	"scyllaDbAssignment/internal/cloud"
	"scyllaDbAssignment/pkg"
)

// ListParams defines the filter criteria
type ListParams struct {
	LT string // less than
	GT string // greater than
}

// Run returns a filtered, tab-separated list of versions
func Run(scyllaClient cloud.ClientInterface, params ListParams) ([]pkg.Version, error) {
	if scyllaClient == nil {
		return nil, errors.New("scylla client is nil")
	}
	versions, err := (scyllaClient).ListVersions()
	if err != nil {
		return nil, err
	}
	filtered := filterVersions(versions, params.GT, params.LT)

	return filtered, nil
}

func filterVersions(versions []pkg.Version, gt, lt string) []pkg.Version {
	var result []pkg.Version
	for _, v := range versions {
		if gt != "" && CompareVersions(v.Name, gt) <= 0 {
			continue
		}
		if lt != "" && CompareVersions(v.Name, lt) >= 0 {
			continue
		}
		result = append(result, v)
	}
	return result
}
