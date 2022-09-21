package git

import (
	// "net/url"
	"path/filepath"

	url "github.com/whilp/git-urls"
)

// RepoURLToPath converts a git repo URL into a workspace path.
func RepoURLToPath(repoURL string) (string, error) {
	// parse URL first
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(append(filepath.SplitList(u.Host), filepath.SplitList(u.Path)...)...)
	return dir, nil
}
