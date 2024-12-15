package status

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	ErrCodeInternalError = iota + 1
	ErrCodeInvalidInput
	ErrCodePermissionDenied
	ErrCodeNotFound
)

// Error struct to hold the error message and code
type Error struct {
	Message string `json:"message"`
	Code    uint   `json:"code"`
}

// Info struct to hold the information message
type Info struct {
	Message string `json:"message"`
}

// NewError creates a new error with a message and code
func NewError(message string, code uint) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

// NewInfo creates a new error with a message and code
func NewInfo(message string) *Info {
	return &Info{
		Message: message,
	}
}

// PrintJsonError method to print error in JSON format and terminate program
func (e *Error) PrintJsonError() {
	errorJSON, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling error: %v", err)
	}
	log.Fatalf("%s", string(errorJSON))
}

// PrintJsonInfo method to print info in JSON format
func (i *Info) PrintJsonInfo(isJsonFormatted bool) {
	if isJsonFormatted {
		fmt.Printf("{\n message: %v \n}", string(i.Message))
	} else {
		infoJSON, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling info: %v", err)
		}
		fmt.Printf("%s\n", string(infoJSON))
	}
}
