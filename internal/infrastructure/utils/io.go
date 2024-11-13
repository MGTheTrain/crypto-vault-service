package utils

import "os"

// File Operations
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func WriteFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}
