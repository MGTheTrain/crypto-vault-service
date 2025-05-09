package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper function to create test files
func CreateTestFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	return nil
}

// Helper function to create a test file and form
func CreateTestFileAndForm(t *testing.T, fileName string, fileContent []byte) (*multipart.Form, error) {
	err := CreateTestFile(fileName, fileContent)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Logf("failed to remove temporary file %s: %v", fileName, err)
		}
	})

	form, err := CreateForm(fileContent, fileName)
	require.NoError(t, err)

	return form, nil
}
