package git

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRemoteBranchNotFound = errors.New("remote branch not found")
)

func ParseCloneError(stderr string, branch string, upstream string) error {
	msg := fmt.Sprintf("fatal: Remote branch %v not found in upstream %v", branch, upstream)
	if strings.Index(string(stderr), msg) >= 0 {
		return ErrRemoteBranchNotFound
	}
	return nil
}
