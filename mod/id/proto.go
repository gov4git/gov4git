package id

import (
	"github.com/gov4git/gov4git/mod"
)

var (
	PublicNS  = mod.RootNS.Sub("id")
	PrivateNS = mod.RootNS.Sub("id")
)

var (
	PublicCredentialsNS  = PublicNS.Sub("public_credentials.json")
	PrivateCredentialsNS = PrivateNS.Sub("private_credentials.json")
)
