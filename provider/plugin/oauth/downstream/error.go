package downstream

import "errors"

// ErrInvalidBearerFormat returned when the header of bearer token is invalid
var ErrInvalidBearerFormat = errors.New("Invalid bearer token format")

// ErrBadRequest returned when client give bad requets
var ErrBadRequest = errors.New("Bad request given by the client")

// ErrInvalidRequest is given when client give unexpected request wanted by the server
var ErrInvalidRequest = errors.New("Invalid request given by the client")

// ErrUnavailable returned when the service is unavailable
var ErrUnavailable = errors.New("Service are temporary unavailable")

// ErrApplicationNotExists returned when client uid or client secret of an application is invalid
var ErrApplicationNotExists = errors.New("Client uid or client secret is invalid")
