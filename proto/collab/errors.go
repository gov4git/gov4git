package collab

import "errors"

var (
	ErrIssueAlreadyExists = errors.New("issue already exists")
	ErrIssueAlreadyClosed = errors.New("issue already closed")
)
