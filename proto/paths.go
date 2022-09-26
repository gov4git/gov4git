package proto

import "path/filepath"

const (
	RootPath = ".ana"
)

var (
	PrivateCredentialsPath = filepath.Join(RootPath, "private_credentials")
)
