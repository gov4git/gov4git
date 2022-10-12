package govproto

import "path/filepath"

type GovUserInfo struct {
	PublicURL string `json:"public_url"`
}

var (
	GovUsersDir         = filepath.Join(GovRoot, "users")
	GovUserInfoFilebase = "info"
	GovUserMetaDirbase  = "meta"
)

func UserInfoFilepath(user string) string {
	return filepath.Join(GovUsersDir, user, GovUserInfoFilebase)
}
