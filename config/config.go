package config

import "path/filepath"

const (
	MainBranch = "main"
	AppPath    = ".gov"
)

var (
	PrivateSoulInfoPath = filepath.Join(AppPath, "private_soul_info")
)
