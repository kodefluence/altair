package entity

import "fmt"

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (eo ErrorObject) Error() string {
	return fmt.Sprintf("Error: %s, Code: %s", eo.Message, eo.Code)
}

type Error struct {
	HttpStatus int
	Errors     []ErrorObject
}

func (e Error) Error() string {
	errorString := "Service error because of:\n"

	for _, err := range e.Errors {
		errorString = errorString + err.Error() + "\n"
	}

	return errorString
}
