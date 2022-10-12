package identityproto

import "path/filepath"

// identity repo paths

var (
	IdentityRoot           = ".gov"
	PublicCredentialsPath  = filepath.Join(IdentityRoot, "public_credentials")
	PrivateCredentialsPath = filepath.Join(IdentityRoot, "private_credentials")
)
