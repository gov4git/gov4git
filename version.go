// Package gov4git provides protocol versioning.
package gov4git

import (
	_ "embed"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var (
	lk          sync.Mutex
	versionInfo = VersionInfo{DirtyBuild: true}
)

func GetVersionInfo() VersionInfo {
	lk.Lock()
	defer lk.Unlock()
	return versionInfo
}

type VersionInfo struct {
	Version    string    `json:"version"`  // Version is the Go version of the protocol module.
	Revision   string    `json:"revision"` // Revision is the git commit hash of the protocol source.
	LastCommit time.Time `json:"last_commit"`
	DirtyBuild bool      `json:"dirty_build"`
}

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	lk.Lock()
	defer lk.Unlock()
	versionInfo.Version = info.Main.Version
	if versionInfo.Version == "(devel)" {
		versionInfo.Version = "dev"
	}
	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			versionInfo.Revision = kv.Value
		case "vcs.time":
			versionInfo.LastCommit, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			versionInfo.DirtyBuild = kv.Value == "true"
		}
	}
}

// Short provides a short string summarizing available version information.
func Short() string {
	vi := GetVersionInfo()
	parts := make([]string, 0, 3)

	if vi.Version != "" {
		parts = append(parts, vi.Version)
	}
	if vi.Revision != "" {
		parts = append(parts, "rev")
		commit := vi.Revision
		if len(commit) > 7 {
			commit = commit[:7]
		}
		parts = append(parts, commit)
		parts = append(parts, vi.LastCommit.Format(time.DateOnly))
	}
	if vi.DirtyBuild {
		parts = append(parts, "dirty")
	}
	if len(parts) == 0 {
		return "dev"
	}
	return strings.Join(parts, "-")
}
