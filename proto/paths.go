package proto

import "path/filepath"

const (
	RootPath           = ".gov"
	LocalAgentPath     = ".gov"
	LocalAgentTempPath = "gov"
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
	GovUserMetaDirbase  = "meta"

	GovGroupsDir         = filepath.Join(GovRoot, "groups")
	GovGroupInfoFilebase = "info"
	GovGroupMetaDirbase  = "meta"

	GovMembersDirbase = "members"

	GovDirPolicyFilebase = "policy"

	GovPollAdFilebase   = "poll_ad"
	GovPollBranchPrefix = "poll"
)
