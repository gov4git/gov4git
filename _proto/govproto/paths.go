package govproto

import (
	"path/filepath"

	"github.com/gov4git/gov4git/proto"
)

// governance repo paths

const (
	GovRoot = ".gov"
)

var (
	GovDirPolicyFilebase = "policy"
)

func SnapshotDir(addr proto.Address, commit string) string {
	return filepath.Join(GovRoot, "snapshot", string(addr.Repo), string(addr.Branch), commit)
}
