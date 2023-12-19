package github

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gov4git/lib4git/must"
)

type Repo struct {
	Owner string `json:"github_repo_owner"`
	Name  string `json:"github_repo_name"`
}

func (x Repo) HTTPS() string {
	return `https://github.com/` + x.Owner + `/` + x.Name
}

// ParseRepo parses a github "owner/repo" pair.
func ParseRepo(ctx context.Context, s string) Repo {
	first, second, ok := strings.Cut(s, "/")
	must.Assertf(ctx, ok, "not a github repo: %v", s)
	return Repo{Owner: first, Name: second}
}

// url can be an HTTPS or an SSH git URL.
func ParseGithubRepoURL(url string) (Repo, error) {
	repo, err := parseGithubRepoHTTPSURL(url)
	if err == nil {
		return repo, nil
	}
	repo, err = parseGithubRepoSSHURL(url)
	if err == nil {
		return repo, nil
	}
	return Repo{}, fmt.Errorf("not an https or ssh github repo url")
}

func parseGithubRepoHTTPSURL(s string) (repo Repo, err error) {
	m := githubRepoHTTPSURLRegexp.FindStringSubmatch(s)
	if m == nil {
		return Repo{}, fmt.Errorf("not a github https repo url")
	}
	return Repo{
		Owner: m[1],
		Name:  m[2],
	}, nil
}

const (
	// this regexp matches url like `https://github.com/gov4git/gov4git.git`
	githubRepoHTTPSURLRegexpSrc = `^https://github\.com/` +
		githubIDRegexp +
		`/` +
		githubIDRegexp +
		`\.git$`
)

var githubRepoHTTPSURLRegexp = regexp.MustCompile(githubRepoHTTPSURLRegexpSrc)

func parseGithubRepoSSHURL(s string) (repo Repo, err error) {
	m := githubRepoSSHURLRegexp.FindStringSubmatch(s)
	if m == nil {
		return Repo{}, fmt.Errorf("not a github ssh repo url")
	}
	return Repo{
		Owner: m[1],
		Name:  m[2],
	}, nil
}

const (
	githubIDRegexp = `([a-zA-Z0-9\.\-_]+)`

	// this regexp matches url like `git@github.com:gov4git/gov4git.git`
	githubRepoSSHURLRegexpSrc = `^git@github\.com:` +
		githubIDRegexp +
		`/` +
		githubIDRegexp +
		`\.git$`
)

var githubRepoSSHURLRegexp = regexp.MustCompile(githubRepoSSHURLRegexpSrc)
