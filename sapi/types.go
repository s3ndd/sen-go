package sapi

// StatusCode is an interface for status code
type StatusCode interface {
	StatusCode() int
}

// ErrorCode is an interface for error code
type ErrorCode interface {
	ErrorCode() int
}

// ErrorField is a map of error fields
type ErrorField map[string]string

// ErrorFields is an interface for error fields
type ErrorFields interface {
	ErrorFields() ErrorField
}

// APIResponseError is an error for API response
type APIResponseError struct {
	statusCode int
	errorCode  int
	err        error
}

// PrivateError is an error for private
type PrivateError struct {
	statusCode int
	message    string
	errorCode  int
	err        error
}

// ValidationError is an error for validation
type ValidationError struct {
	statusCode int
	errorCode  int
	message    string
	fields     ErrorField
}

// ErrorResponse is a response for error
type ErrorResponse struct {
	Error  string     `json:"error"`
	Code   int        `json:"code,omitempty"`
	Fields ErrorField `json:"fields,omitempty"`
}

// RequestForLog is a request for log
type RequestForLog struct {
	RequestID string
	Method    string
	URL       string
	Host      string
	ClientIP  string
}
