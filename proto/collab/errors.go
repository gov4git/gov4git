package collab

import "errors"

var (
	ErrConcernAlreadyExists = errors.New("concern already exists")
	ErrConcernAlreadyClosed = errors.New("concern already closed")
)
