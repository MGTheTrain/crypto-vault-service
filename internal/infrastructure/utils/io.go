package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

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

// Helper function to create multipart form with multiple files
func CreateMultipleFilesForm(contents [][]byte, fileNames []string) (*multipart.Form, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Ensure the number of contents and file names match
	if len(contents) != len(fileNames) {
		return nil, fmt.Errorf("number of contents and file names must match")
	}

	// Loop over the contents and write each one to the multipart form
	for i, content := range contents {
		fileWriter, err := writer.CreateFormFile("files", fileNames[i])
		if err != nil {
			return nil, err
		}

		_, err = fileWriter.Write(content)
		if err != nil {
			return nil, err
		}
	}

	// Close the writer
	writer.Close()

	// Read the form from the buffer
	mr := multipart.NewReader(&buf, writer.Boundary())
	form, err := mr.ReadForm(10 << 20)
	if err != nil {
		return nil, err
	}

	return form, nil
}
