package govproto

import "path/filepath"

// governance repo paths

const (
	GovRoot = ".gov"
)

var (
	GovDirPolicyFilebase = "policy"
)

func SnapshotDir(repo string, commit string) string {
	return filepath.Join(GovRoot, "snapshot", repo, commit)
}
