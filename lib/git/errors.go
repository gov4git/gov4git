package git

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRemoteBranchNotFound = errors.New("remote branch not found")
	ErrNothingToCommit      = errors.New("nothing to commit")
	ErrAlreadyOnBranch      = errors.New("already on branch")
)

func ParseCloneError(stderr string, branch string, upstream string) error {
	msg := fmt.Sprintf("fatal: Remote branch %v not found in upstream %v", branch, upstream)
	if strings.Index(string(stderr), msg) >= 0 {
		return ErrRemoteBranchNotFound
	}
	return nil
}

func ParseCommitError(stderr string) error {
	if strings.Index(stderr, "nothing to commit, working tree clean") >= 0 {
		return ErrNothingToCommit
	}
	return nil
}

func ParseCheckoutError(stderr string) error {
	if strings.Index(stderr, "Already on") >= 0 {
		return ErrAlreadyOnBranch
	}
	return nil
}
