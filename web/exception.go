package web

import (
	"errors"
	"fmt"
	"net/http"
)

// HTTPError http error struct
type HTTPError struct {
	Error   error
	Message string
	Code    int
}

func (h *HTTPError) String() string {
	return fmt.Sprintf("%d - %s", h.Code, h.Message)
}

type httpHandler func(ctx *Context) *HTTPError

func newHTTPError(err error, msg string, code int) *HTTPError {
	return &HTTPError{err, msg, code}
}

func serverError(err error) *HTTPError {
	return newHTTPError(err, err.Error(), http.StatusInternalServerError)
}

func badRequestError(msg string) *HTTPError {
	err := errors.New(msg)
	return newHTTPError(err, err.Error(), http.StatusBadRequest)
}

func notFoundError(err error, msg string) *HTTPError {
	return newHTTPError(err, msg, http.StatusNotFound)
}
