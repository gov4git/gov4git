package proto

import "path/filepath"

const (
	LocalAgentPath     = ".gov4git"
	LocalAgentTempPath = "gov4git"
)

// identity repo paths

var (
	IdentityRoot           = ".gov"
	PublicCredentialsPath  = filepath.Join(IdentityRoot, "public_credentials")
	PrivateCredentialsPath = filepath.Join(IdentityRoot, "private_credentials")
)

// governance repo paths

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

func UserInfoFilepath(user string) string {
	return filepath.Join(GovUsersDir, user, GovUserInfoFilebase)
}

func GroupInfoFilepath(group string) string {
	return filepath.Join(GovGroupsDir, group, GovGroupInfoFilebase)
}

func GroupMemberFilepath(group string, user string) string {
	return filepath.Join(GovGroupsDir, group, GovMembersDirbase, user)
}

func SnapshotDir(repo string, commit string) string {
	return filepath.Join(GovRoot, "snapshot", repo, commit)
}
