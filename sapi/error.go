package sapi

import "net/http"

// NewAPIResponseError creates a new APIResponseError
func NewAPIResponseError(err error, statusCode int) APIResponseError {
	return APIResponseError{
		statusCode: statusCode,
		err:        err,
	}
}

// WithErrorCode sets the error code for the APIResponseError
func (e APIResponseError) WithErrorCode(code int) APIResponseError {
	e.errorCode = code
	return e
}

// ErrorCode returns the error code for the APIResponseError
func (e APIResponseError) Error() string {
	return e.err.Error()
}

// ErrorCode returns the error code for the APIResponseError
func (e APIResponseError) ErrorCode() int {
	return e.errorCode
}

// StatusCode returns the status code for the APIResponseError
func (e APIResponseError) StatusCode() int {
	return e.statusCode
}

// NewPrivateError creates a new PrivateError
func NewPrivateError(err error) PrivateError {
	return NewPrivateErrorWithStatusCode(err, http.StatusInternalServerError)
}

// NewPrivateErrorWithStatusCode creates a new PrivateError with a status code
func NewPrivateErrorWithStatusCode(err error, statusCode int) PrivateError {
	return PrivateError{
		statusCode: statusCode,
		err:        err,
	}
}

// WithMessage sets the message for the PrivateError
func (e PrivateError) WithMessage(message string) PrivateError {
	e.message = message
	return e
}

// WithErrorCode sets the error code for the PrivateError
func (e PrivateError) WithErrorCode(code int) PrivateError {
	e.errorCode = code
	return e
}

// ErrorCode returns the error code for the PrivateError
func (e PrivateError) ErrorCode() int {
	return e.errorCode
}

// StatusCode returns the status code for the PrivateError
func (e PrivateError) Error() string {
	return e.err.Error()
}

// StatusCode returns the status code for the PrivateError
func (e PrivateError) StatusCode() int {
	return e.statusCode
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, statusCode int) ValidationError {
	if statusCode == 0 {
		statusCode = http.StatusUnprocessableEntity
	}

	return ValidationError{
		statusCode: statusCode,
		message:    message,
		fields:     ErrorField{},
	}
}

// WithErrorCode sets the error code for the ValidationError
func (e ValidationError) WithErrorCode(code int) ValidationError {
	e.errorCode = code
	return e
}

// WithErrorFields sets the error fields for the ValidationError
func (e ValidationError) WithErrorFields(errorFields ErrorField) ValidationError {
	e.fields = errorFields
	return e
}

// ErrorCode returns the error code for the ValidationError
func (e ValidationError) ErrorCode() int {
	return e.errorCode
}

// StatusCode returns the status code for the ValidationError
func (e ValidationError) StatusCode() int {
	return e.statusCode
}

// Error returns the error message for the ValidationError
func (e ValidationError) Error() string {
	return e.message
}

// ErrorFields returns the error fields for the ValidationError
func (e ValidationError) ErrorFields() ErrorField {
	return e.fields
}
