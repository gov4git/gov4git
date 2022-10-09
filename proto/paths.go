package proto

import "path/filepath"

const (
	LocalAgentPath     = ".gov4git"
	LocalAgentTempPath = "gov4git"
)

var (
	IdentityRoot           = ".gov"
	PublicCredentialsPath  = filepath.Join(IdentityRoot, "public_credentials")
	PrivateCredentialsPath = filepath.Join(IdentityRoot, "private_credentials")
)

// governance-related constants

const (
	GovRoot = ".gov"
)

var (
	GovUsersDir         = filepath.Join(GovRoot, "users")
	GovUserInfoFilebase = "info"
	GovUserMetaDirbase  = "meta"

	GovGroupsDir         = filepath.Join(GovRoot, "groups")
	GovGroupInfoFilebase = "info"
	GovGroupMetaDirbase  = "meta"

	GovMembersDirbase = "members"

	GovDirPolicyFilebase = "policy"
)

func SnapshotDir(repo string, commit string) string {
	return filepath.Join(GovRoot, "snapshot", repo, commit)
}
