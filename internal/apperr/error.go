package apperr

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Doris-Mwito5/ginja-ai/internal/logger"
	"net/http"
)

type Type string

// error types
const (
	Authorization        Type = "AUTHORIZATION" // Authentication Failures -
	BadRequest           Type = "BAD_REQUEST"   // Validation errors / BadInput
	Conflict             Type = "CONFLICT"      // Already exists (eg, create account with existent email) - 409
	Internal             Type = "INTERNAL"      // Server (500) and fallback errors
	Permission           Type = "PERMISSION_DENIED"
	NotFound             Type = "NOT_FOUND"              // For not finding resource
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"      // for uploading tons of JSON, or an image over the limit - 413
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"    // For long running handlers
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE" // for http 415
	UnexpextedError      Type = "UNEXPECTED_ERROR"
	DatabaseError        Type = "DATABASE_ERROR"
)

var (
	ErrNetwork               = errors.New("network request failed")
	ErrConnectionRefused     = errors.New("connection refused")
	ErrConnectionResetByPeer = errors.New("connection reset by peer")
	ErrDialTcpIOTimeout      = errors.New("dial tcp io timeout")
)

type Error struct {
	Type        Type     `json:"type"`
	Service     string   `json:"service"`
	Message     string   `json:"message"`
	RequestID   string   `json:"request_id"`
	IPAddress   string   `json:"ip_address"`
	MACAddress  string   `json:"mac_address"`
	Payload     any      `json:"payload"`
	LogMessages []string `json:"log_messages"`
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case UnexpextedError:
		return http.StatusExpectationFailed
	case DatabaseError:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// return regular err
func (e *Error) Error() string {
	return e.Message
}

// message
func (e *Error) GetMessage() string {
	return e.Message
}

// NewAuthorization to create a 401
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: reason,
	}
}

// NewConflict to create an error for 409
func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal(message string) *Error {
	return &Error{
		Type:    Internal,
		Message: message,
	}
}

func NewPermission(reason string) *Error {
	return &Error{
		Type:    Permission,
		Message: reason,
	}
}

// NewNotFound to create an error for 404
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("%v not found", name),
	}
}

// NewPayloadTooLarge to create an error for 413
func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

// NewServiceUnavailable to create an error for 503
func NewServiceUnavailable(name string) *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: fmt.Sprintf("service %s unavailable or timed out", name),
	}
}

// NewUnsupportedMediaType to create an error for 415
func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

func NewUnexpectedError(reason string) *Error {
	return &Error{
		Type:    UnexpextedError,
		Message: reason,
	}
}

func NewDatabaseError(err error) *Error {
	appError := CastError(err, Internal)
	if IsNoRowsErr(err) {
		appError.Type = NotFound
		appError.Message = err.Error()
	} else {
		appError.Type = Internal
		appError.Message = err.Error()
	}

	return appError
}

func NewError(err error) *Error {
	appError, ok := err.(*Error)
	if !ok {
		return NewErrorWithType(
			err,
			Internal,
		)
	}

	return appError
}

func NewErrorWithType(
	err error,
	errorType Type,
) *Error {

	if errorType == "" {
		errorType = Internal
	}

	appError, ok := err.(*Error)
	if !ok {
		appError = &Error{
			Message: err.Error(),
			Type:    errorType,
		}
	}

	appError.Type = errorType

	return appError
}

func CastError(err error, errorType Type) *Error {
	return New(err, errorType)
}

func New(err error, errorType Type) *Error {
	if errorType == "" {
		errorType = UnexpextedError
	}

	appError, ok := err.(*Error)
	if !ok {
		appError = &Error{
			Message: err.Error(),
			Type:    errorType,
		}
	}

	return appError
}

func IsNoRowsErr(err error) bool {
	appError := CastError(err, Internal)
	return appError.Error() == sql.ErrNoRows.Error()
}

func (e *Error) LogErrorMessage(message string, values ...interface{}) error {
	logger.Errorf(message, values...)
	return e
}

func LogErrorMessage(service, message string, err error) {
	logger.Errorf(service, message, err)
}

func (e *Error) JsonResponse() map[string]interface{} {

	if e.Message == "" {
		e.Message = "Failed to perform request. Please try again."
	}

	jsonResponse := map[string]interface{}{
		"error_code":    e.Type,
		"error_message": e.Message,
		"status_code":   e.Status(),
	}

	return jsonResponse
}
