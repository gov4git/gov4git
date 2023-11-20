package github

import "github.com/google/go-github/v55/github"

func IsLabelAlreadyExists(err error) bool {

	gerr, ok := err.(*github.ErrorResponse)
	if !ok {
		return false
	}
	if len(gerr.Errors) != 1 {
		return false
	}
	if gerr.Errors[0].Code != "already_exists" {
		return false
	}
	return true
}
