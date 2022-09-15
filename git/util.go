package git

import (
	// "net/url"
	"path/filepath"

	url "github.com/whilp/git-urls"
)

// URL2Path converts a git repo URL into a workspace path.
func URL2Path(repoURL string) (string, error) {
	// parse URL first
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(append(filepath.SplitList(u.Host), filepath.SplitList(u.Path)...)...)
	return dir, nil
}
