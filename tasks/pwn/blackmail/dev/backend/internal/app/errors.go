// Domain errors
package app

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrExists     = errors.New("already exists")
	ErrForbidden  = errors.New("forbidden")
	ErrWrongCreds = errors.New("wrong username or password")
)
