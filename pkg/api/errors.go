package api

import "errors"

var (
	// generic
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrUnknown   = errors.New("unknown")

	// auth
	ErrLoginFailed   = errors.New("login failed")
	ErrLoginRequired = errors.New("login required")
	ErrUnauthorized  = errors.New("unauthorized")

	// parsing
	ErrParsingFailed = errors.New("parsing failed")

	// service
	ErrServiceUnavailable = errors.New("service unavailable")

	// request
	ErrRequestFailed        = errors.New("request failed")
	ErrTimeout              = errors.New("timeout")
	ErrNoAuthHeaderIncluded = errors.New("no auth header included")
	ErrMalformedAuthHeader  = errors.New("malformed authorization header")

	// config
	ErrReadConfig    = errors.New("failed to read config")
	ErrInvalidConfig = errors.New("invalid config")

	// proxy
	ErrServiceNotFound = errors.New("service not found")
	ErrInvalidPath     = errors.New("invalid path")
	ErrInvalidMethod   = errors.New("invalid method")
	ErrInvalidRequest  = errors.New("invalid request")
)
