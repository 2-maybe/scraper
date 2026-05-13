package urlsanitizer

import (
	"strings"
	"testing"
)

func TestArg(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		errContains string
	}{
		{
			name:     "extract path",
			input:    "https://www.google.com/",
			expected: "www.google.com",
		},
		{
			name:     "transform extract uppercase paths",
			input:    "https://www.GOOGLE.com/",
			expected: "www.google.com",
		},
		{
			name:        "Block http",
			input:       "http://httpforever.com/",
			errContains: "Http Path not allowed",
		},
		{
			name:        "Invalid Url",
			input:       "https://limozine.netkw",
			errContains: "Malfomred Url structure",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := NormalizeUrl(tc.input)

			if tc.errContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, but got nil", tc.errContains)
				}
				if !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("expected error %q, but got %q", tc.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if data != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, data)
			}
		})
	}
}
