package app_errors

import "errors"

var (
	// Service
	ErrIdNotProvided      = errors.New("There is not id. Provide one please!")
	ErrDataIsNil          = errors.New("Malformed data. Please fix the input data provided")
	ErrInvalidCredentials = errors.New("Invalid credentials!")
)
