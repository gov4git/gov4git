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

func TestParseGithubRepoSSHURL(t *testing.T) {
	repo, err := parseGithubRepoSSHURL("git@github.com:abc/x.y.z.git")
	if err != nil {
		t.Error(err)
	}
	if repo.Owner != "abc" {
		t.Errorf("expecting %v, got %v", "abc", repo.Owner)
	}
	if repo.Name != "x.y.z" {
		t.Errorf("expecting %v, got %v", "x.y.z", repo.Name)
	}
}
