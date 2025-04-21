package events

import "errors"

var (
	// ErrDriverNotFound is returned when a requested driver is not found
	ErrDriverNotFound = errors.New("driver not found")
	// ErrDriverNotInitialized is returned when a driver is used before initialization
	ErrDriverNotInitialized = errors.New("driver not initialized")
	// ErrNoDefaultDriver is returned when no default driver is set
	ErrNoDefaultDriver = errors.New("no default driver set")
)
