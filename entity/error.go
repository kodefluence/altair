package entity

import "fmt"

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (re ResponseError) Error() string {
	return fmt.Sprintf("Error: %s, Code: %s", re.Message, re.Code)
}

type Error struct {
	HttpStatus int
	Errors     []ResponseError
}

func (e Error) Error() string {
	errorString := "Service error because of:\n"

	for _, err := range e.Errors {
		errorString = errorString + err.Error() + "\n"
	}

	return errorString
}
