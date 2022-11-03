package govproto

import "path/filepath"

type GovGroupInfo struct {
}

var (
	GovGroupsDir         = filepath.Join(GovRoot, "groups")
	GovGroupInfoFilebase = "info"
	GovGroupMetaDirbase  = "meta"
	GovMembersDirbase    = "members"
)

func GroupInfoFilepath(group string) string {
	return filepath.Join(GovGroupsDir, group, GovGroupInfoFilebase)
}

func GroupMemberFilepath(group string, user string) string {
	return filepath.Join(GovGroupsDir, group, GovMembersDirbase, user)
}
