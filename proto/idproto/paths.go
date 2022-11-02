package idproto

import "path/filepath"

var (
// IdentityBranch = "main" // XXX: for identity public and private repos
)

// identity repo paths

var (
	IdentityRoot           = ".id"
	PublicCredentialsPath  = filepath.Join(IdentityRoot, "public_credentials.json")
	PrivateCredentialsPath = filepath.Join(IdentityRoot, "private_credentials.json")
)
