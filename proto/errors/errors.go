package errors

import "errors"

var (
	ErrGithubUnreachable           = errors.New("github is unreachable")
	ErrCommunityPublicUnreachable  = errors.New("community public repo is unreachable")
	ErrCommunityPrivateUnreachable = errors.New("community private repo is unreachable")
	ErrMemberPublicUnreachable     = errors.New("member public repo is unreachable")
	ErrMemberPrivateUnreachable    = errors.New("member private repo is unreachable")
)
