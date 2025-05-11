package testutils

import (
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// CreateTestFile create a test files
func CreateTestFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0600)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	return nil
}

// CreateTestFileAndForm creates a test file and form
func CreateTestFileAndForm(t *testing.T, fileName string, fileContent []byte) (*multipart.Form, error) {
	err := CreateTestFile(fileName, fileContent)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Logf("failed to remove temporary file %s: %v", fileName, err)
		}
	})

	form, err := utils.CreateForm(fileContent, fileName)
	require.NoError(t, err)

	return form, nil
}
