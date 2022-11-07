// Package member implements community member management services
package member

import (
	"github.com/gov4git/gov4git/mod"
)

const (
	everybody = "everybody"
)

var (
	userNS  = mod.NS("users")
	groupNS = mod.NS("groups")
)
