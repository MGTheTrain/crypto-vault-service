package utils

import (
	"bytes"
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
