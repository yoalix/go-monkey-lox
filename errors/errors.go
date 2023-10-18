package errors

import (
	"errors"
	"fmt"
)

var HadError bool

func Error(line int, message string) {
	Report(line, "", message)
}

// TODO: add row and column number and where in the line the error is
func Report(line int, where string, message string) string {
	errorMessage := fmt.Sprintf("[line %d] Error %s: %s\n", line, where, message)
	err := errors.New(errorMessage)
	fmt.Println(err)
	HadError = true
	return errorMessage
}
