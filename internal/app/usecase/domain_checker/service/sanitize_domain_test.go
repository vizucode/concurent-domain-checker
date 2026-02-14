package service

import (
	"reflect"
	"sort"
	"testing"
)

func TestSanitizeDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Basic Sanitation",
			input:    []string{"Example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Remove Protocol HTTP",
			input:    []string{"http://example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Remove Protocol HTTPS",
			input:    []string{"https://example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Remove WWW",
			input:    []string{"www.example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Remove WWW with Protocol",
			input:    []string{"https://www.example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Trim Whitespace",
			input:    []string{"  example.com  "},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Empty String",
			input:    []string{""},
			expected: []string{},
		},
		{
			name:     "Duplicate Removal",
			input:    []string{"example.com", "example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Mixed Case Duplicate",
			input:    []string{"example.com", "EXAMPLE.COM"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Multiple Valid Domains",
			input:    []string{"example.com", "google.com"},
			expected: []string{"https://example.com", "https://google.com"},
		},
		{
			name:     "Complex Case",
			input:    []string{"  HtTp://WwW.ExAmPlE.cOm  "},
			expected: []string{"https://example.com"},
		},
		{
			name:     "Mixed Empty and Valid",
			input:    []string{"", "example.com", "   "},
			expected: []string{"https://example.com"},
		},
	}

	s := &domainCheckerService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultChan := s.sanitizeDomain(tt.input)

			var got []string
			for domain := range resultChan {
				got = append(got, domain)
			}

			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}

			// Sort to ensure deterministic comparison
			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("sanitizeDomain() = %v, want %v", got, tt.expected)
			}
		})
	}
}
