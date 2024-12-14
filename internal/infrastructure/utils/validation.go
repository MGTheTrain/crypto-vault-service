package utils

import "fmt"

// CheckNonEmptyStrings checks if any of the input strings are empty
func CheckNonEmptyStrings(args ...string) error {
	for _, arg := range args {
		if arg == "" {
			return fmt.Errorf("one of the input strings is empty")
		}
	}
	return nil
}
