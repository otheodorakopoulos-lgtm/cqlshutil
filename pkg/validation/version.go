package validation

import (
	"fmt"
	"regexp"
	"strconv"
)

// ValidateVersion checks if a version string matches the expected patterns.
// Accepted forms:
//
//	X
//	X.Y
//	X.Y.Z
//	X.Y.Z~rcN
func ValidateVersion(v string) error {
	pattern := `^(\d+)(\.\d+){0,2}(~rc\d+)?$`

	re := regexp.MustCompile(pattern)
	if !re.MatchString(v) {
		return errInvalidVersion(v)
	}
	return nil
}

func ValidateFullVersion(input string) error {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:~rc(\d+))?$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return fmt.Errorf("invalid version format. Version is not full version: %s", input)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	if major <= 0 || minor < 0 || patch < 0 {
		return fmt.Errorf("version numbers must be > 0: %s", input)
	}

	if matches[4] != "" {
		n, _ := strconv.Atoi(matches[4])
		if n < 0 {
			return fmt.Errorf("rc number must be > 0: %s", input)
		}
	}

	return nil
}

func errInvalidVersion(input string) error {
	return fmt.Errorf("invalid version format: %q (expected X, X.Y, X.Y.Z, or X.Y.Z~rcN)", input)
}
