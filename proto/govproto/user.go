package govproto

import (
	"path/filepath"

	"github.com/gov4git/gov4git/proto"
)

type GovUserInfo struct {
	Address proto.Address `json:"address"`
}

var (
	GovUsersDir         = filepath.Join(GovRoot, "users")
	GovUserInfoFilebase = "info"
	GovUserMetaDirbase  = "meta"
)

func UserInfoFilepath(user string) string {
	return filepath.Join(GovUsersDir, user, GovUserInfoFilebase)
}
