package listing

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"2025.1.0", "2025.1.0", 0},
		{"2025.1.0~rc1", "2025.1.0~rc1", 0},

		{"2025.1.1", "2025.1.0", 1},
		{"2025.1.0", "2025.1.1", -1},
		{"2025.2.0", "2025.1.9", 1},

		{"2025.1.0", "2025.1.0~rc1", 1},
		{"2025.1.0~rc1", "2025.1.0", -1},

		{"2025.1.0~rc1", "2025.1.0~rc2", -1},
		{"2025.1.0~rc2", "2025.1.0~rc1", 1},

		{"2025.1", "2025.1.0", 0},
		{"2025.1", "2025.1.1", -1},
		{"2025.1.2", "2025.1", 1},

		{"2025.1.0~rc10", "2025.1.0~rc2", 1},
		{"2025.1.0~rc2", "2025.1.0~rc10", -1},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			got := CompareVersions(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %d; want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}
