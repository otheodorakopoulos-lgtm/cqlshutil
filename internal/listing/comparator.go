package listing

import (
	"strconv"
	"strings"
)

func CompareVersions(a, b string) int {
	aParts, aRC := splitVersion(a)
	bParts, bRC := splitVersion(b)

	for i := 0; i < maximumOf(len(aParts), len(bParts)); i++ {
		var ai, bi int
		if i < len(aParts) {
			ai = aParts[i]
		}
		if i < len(bParts) {
			bi = bParts[i]
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}

	if aRC == "" && bRC != "" {
		return 1 // a is stable, b is rc
	}
	if aRC != "" && bRC == "" {
		return -1 // a is rc, b is stable
	}
	if aRC != "" {
		ai, _ := strconv.Atoi(aRC)
		bi, _ := strconv.Atoi(bRC)
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}

	return 0
}

func splitVersion(v string) (nums []int, rc string) {
	if idx := strings.Index(v, "~rc"); idx != -1 {
		rc = v[idx+3:]
		v = v[:idx]
	}
	parts := strings.Split(v, ".")
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		nums = append(nums, n)
	}
	return
}

func maximumOf(a, b int) int {
	if a > b {
		return a
	}
	return b
}
