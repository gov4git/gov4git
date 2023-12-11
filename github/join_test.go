package github

import "testing"

func TestParseGithubRepoHTTPSURL(t *testing.T) {
	repo, err := parseGithubRepoHTTPSURL("https://github.com/abc/xyz.git")
	if err != nil {
		t.Error(err)
	}
	if repo.Owner != "abc" {
		t.Errorf("expecting %v, got %v", "abc", repo.Owner)
	}
	if repo.Name != "xyz" {
		t.Errorf("expecting %v, got %v", "xyz", repo.Name)
	}
}
