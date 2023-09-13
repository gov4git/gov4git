package collab

import "errors"

var (
	ErrMotionAlreadyExists = errors.New("motion already exists")
	ErrMotionAlreadyClosed = errors.New("motion already closed")
)
