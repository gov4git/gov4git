package docket

import "errors"

var (
	ErrMotionAlreadyExists = errors.New("motion already exists")
	ErrMotionAlreadyClosed = errors.New("motion already closed")
	ErrMotionNotClosed     = errors.New("motion is not closed")
	ErrMotionAlreadyFrozen = errors.New("motion already frozen")
	ErrMotionNotFrozen     = errors.New("motion is not frozen")
)
