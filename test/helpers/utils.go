package helpers

import (
	"bytes"
	"fmt"
	"mime/multipart"
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

// Helper function to create multipart form
func CreateForm(content []byte, fileName string) (*multipart.Form, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fileWriter, err := writer.CreateFormFile("files", fileName)
	if err != nil {
		return nil, err
	}

	_, err = fileWriter.Write(content)
	if err != nil {
		return nil, err
	}

	writer.Close()

	mr := multipart.NewReader(&buf, writer.Boundary())
	form, err := mr.ReadForm(10 << 20)
	if err != nil {
		return nil, err
	}

	return form, nil
}
