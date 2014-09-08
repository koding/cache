// Package cache provides basic caching mechanisms for Go(lang) projects.
package cache

import "errors"

var (
	// ErrNotFound holds exported `not found error` for not found items
	ErrNotFound = errors.New("not found")
)
