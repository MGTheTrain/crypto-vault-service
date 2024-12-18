package helpers

import (
	"fmt"
	"os"
)

// Helper function to create test files
func CreateTestFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	return nil
}
