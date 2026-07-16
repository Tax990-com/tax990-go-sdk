// Package utils provides helper utilities for EIN validation and webhook verification.
package utils

import (
	"fmt"
	"regexp"
)

var einPattern = regexp.MustCompile(`^(\d{9}|\d{2}-\d{7})$`)

// ValidateEIN reports whether the given EIN matches the expected format.
// Accepted formats: 9 consecutive digits ("123456789") or formatted ("12-3456789").
func ValidateEIN(ein string) bool {
	return einPattern.MatchString(ein)
}

// FormatEIN converts a 9-digit EIN to the standard XX-XXXXXXX format.
// If the EIN is already formatted or invalid, it is returned unchanged.
func FormatEIN(ein string) string {
	if len(ein) == 9 {
		allDigits := true
		for _, r := range ein {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return fmt.Sprintf("%s-%s", ein[:2], ein[2:])
		}
	}
	return ein
}
