package apierrors

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ouiasy/golang-auth/httputils"
)

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	// todo: get request id, implement middleware generating requestid
	ctx := r.Context()
	var e *HTTPError
	switch {
	case errors.As(err, &e):
		slog.Log(ctx, e.InternalLevel, e.InternalError.Error()) // todo: add request id
		httputils.SendJSON(w, e.Code, e)
	}
}

func BadRequestError(msg ApiError) *HTTPError {
	return httpError(http.StatusBadRequest, msg)
}

func UnprocessableEntityError(msg ApiError) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, msg)
}

func TooManyRequestsError(msg ApiError) *HTTPError {
	return httpError(http.StatusTooManyRequests, msg)
}

func InternalServerError(msg ApiError) *HTTPError {
	return httpError(http.StatusInternalServerError, msg)
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	RequestID string `json:"request_id,omitempty"`

	InternalLevel   slog.Level `json:"-"`
	InternalError   error      `json:"-"`
	InternalMessage string     `json:"-"`
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// WithInternalError adds internal error information to the error
func (e *HTTPError) WithInternalError(level slog.Level, err error) *HTTPError {
	e.InternalError = err
	e.InternalLevel = level
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

func httpError(code int, msg ApiError) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: msg.Message,
	}
}
