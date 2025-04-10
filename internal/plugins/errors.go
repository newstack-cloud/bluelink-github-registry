package plugins

import "errors"

var (
	// ErrUnauthorised is returned when a user is not authorised
	// to access information about a plugin repository.
	ErrUnauthorised = errors.New("not authorised to access this repository")

	// ErrForbidden is returned when a user is forbidden
	// from accessing information about a plugin repository.
	ErrForbidden = errors.New("forbidden to access this repository")

	// ErrRepoNotFound is returned when a plugin repository
	// cannot be found.
	ErrRepoNotFound = errors.New("plugin repository not found")
)
