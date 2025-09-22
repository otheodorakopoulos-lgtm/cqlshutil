package validation

import "testing"

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		// valid versions
		{"2025", false},
		{"2025.1", false},
		{"2025.1.2", false},
		{"2025.1.2~rc1", false},
		{"0.0.0", false},

		// invalid versions
		{"", true},
		{"2025..1", true},
		{"2025.a", true},
		{"2025.1.2~rc", true},
		{"2025.1.2~rc-1", true},
		{"2025.1.2~rcX", true},
		{"2025.1.2~1", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ValidateVersion(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateVersion(%q) error = %v, want error? %v", tt.input, err, tt.expectError)
			}
		})
	}
}

func TestValidateFullVersion(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		// valid full versions
		{"2025.1.0", false},
		{"2025.1.0~rc1", false},
		{"1.0.0", false},
		{"10.5.3~rc2", false},

		// invalid full versions
		{"2025", true},
		{"2025.1", true},
		{"2025.1.0~rc", true},
		{"0.1.2", true},      // major = 0
		{"1.-1.2", true},     // minor < 0
		{"1.1.-2", true},     // patch < 0
		{"1.1.1~rc-1", true}, // rc < 0
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ValidateFullVersion(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateFullVersion(%q) error = %v, want error? %v", tt.input, err, tt.expectError)
			}
		})
	}
}
