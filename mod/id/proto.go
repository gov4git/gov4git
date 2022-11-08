package id

import "github.com/gov4git/gov4git/lib/ns"

var (
	PublicNS  = ns.NS(".id")
	PrivateNS = ns.NS(".id")
)

var (
	PublicCredentialsNS  = PublicNS.Sub("public_credentials.json")
	PrivateCredentialsNS = PrivateNS.Sub("private_credentials.json")
)
