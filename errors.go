// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }

const (
	ErrNotAFile       = Error("not a file")
	ErrNotADirectory  = Error("not a directory")
	ErrNotImplemented = Error("not implemented")
)
