package utils

import "errors"

/* Custom errors for REST API */

// ErrBadRequest is used when the request is malformed
var ErrBadRequest = errors.New("bad request")

// ErrNotFound is used when a requested entity does not exist
var ErrNotFound = errors.New("not found")
