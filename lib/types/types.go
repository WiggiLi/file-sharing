package types

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound 		   = errors.New("resource not found")
	ErrConflict 		   = errors.New("datamodel conflict")
	ErrForbidden           = errors.New("forbidden access")
	ErrNeedMore            = errors.New("need more input")
	ErrBadRequest          = errors.New("bad request")
	ErrPartialOk           = errors.New("partial okay")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrGone                = errors.New("resource gone")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrNotAllowed          = errors.New("operation not allowed")
	ErrBusy                = errors.New("resource is busy")
	ErrUnauthorized        = errors.New("unauthorized")
)

//create custom HTTP error for tests
type HTTPError struct{
	Code int
	Message string
}

// HTTPCode returns the HTTP code of a given custom HTTP error, with 500 as default.
func HTTPCode(err error) int {
	code := 500
	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}
	return code
}