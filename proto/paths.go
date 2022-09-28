package proto

import "path/filepath"

// soul-related constants

const (
	RootPath           = ".ana"
	LocalAgentPath     = ".ana"
	LocalAgentTempPath = "ana"
)

var (
	PublicCredentialsPath  = filepath.Join(RootPath, "public_credentials")
	PrivateCredentialsPath = filepath.Join(RootPath, "private_credentials")
)

// governance-related constants

const (
	GovRoot = ".gov"
)

var (
	GovUsersDir         = filepath.Join(GovRoot, "users")
	GovUserInfoFilebase = "info"
)
