package service

import (
	"strings"
)

func (s *domainCheckerService) sanitizeDomain(domains []string) chan string {
	result := make(chan string)

	go func() {
		defer close(result)
		seen := make(map[string]struct{})

		for _, valueDomain := range domains {
			valueDomain = strings.TrimSpace(valueDomain)
			valueDomain = strings.ToLower(valueDomain)

			valueDomain = strings.TrimPrefix(valueDomain, "http://")
			valueDomain = strings.TrimPrefix(valueDomain, "https://")

			valueDomain = strings.TrimPrefix(valueDomain, "www.")

			if valueDomain == "" {
				continue
			}

			normalized := "https://" + valueDomain

			if _, exists := seen[normalized]; exists {
				continue
			}

			seen[normalized] = struct{}{}
			result <- normalized
		}
	}()

	return result
}
