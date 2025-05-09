package utils

import "strconv"

// Helper function to safely convert a string to an integer with error handling
func ConvertToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0 // Return 0 if the conversion fails, or you can choose a different default value
	}
	return value
}

// Helper function to safely convert a string to an int64 with error handling
func ConvertToInt64(str string) int64 {
	value, err := strconv.ParseInt(str, 10, 64) // Parsing as int64
	if err != nil {
		return 0 // Return 0 if the conversion fails, or you can choose a different default value
	}
	return value
}
