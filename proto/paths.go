package proto

import "path/filepath"

const (
	RootPath = ".ana"
)

var (
	PublicCredentialsPath  = filepath.Join(RootPath, "public_credentials")
	PrivateCredentialsPath = filepath.Join(RootPath, "private_credentials")
)
