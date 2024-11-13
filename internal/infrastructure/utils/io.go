package utils

import (
	"io/ioutil"
)

// File Operations
func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func WriteFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}
