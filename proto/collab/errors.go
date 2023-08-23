package collab

import "errors"

var (
	ErrConcernAlreadyExists = errors.New("concern already exists")
	ErrConcernAlreadyClosed = errors.New("concern already closed")

	ErrProposalAlreadyExists = errors.New("proposal already exists")
	ErrProposalAlreadyClosed = errors.New("proposal already closed")
)
