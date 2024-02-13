package github

import (
	"context"
	"sync"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(ctx context.Context, repo Repo) *github.Client {
	ts := GetTokenSource(ctx, repo)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

//

func MakeStaticTokenSource(ctx context.Context, accessToken string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
}

// token source manager in context

type contextKeyTokenSourceManager struct{}

func WithTokenSource(ctx context.Context, m *TokenSourceManager) context.Context {
	if m == nil {
		m = NewTokenSourceManager()
	}
	return context.WithValue(ctx, contextKeyTokenSourceManager{}, m)
}

func SetTokenSource(ctx context.Context, repo Repo, a oauth2.TokenSource) {
	ctx.Value(contextKeyTokenSourceManager{}).(*TokenSourceManager).SetTokenSource(repo, a)
}

func GetTokenSource(ctx context.Context, repo Repo) oauth2.TokenSource {
	if am, ok := ctx.Value(contextKeyTokenSourceManager{}).(*TokenSourceManager); ok {
		return am.GetTokenSource(repo)
	}
	return nil
}

// TokenSourceManager provides authentication methods given a repo URL.
type TokenSourceManager struct {
	lk  sync.Mutex
	url map[Repo]oauth2.TokenSource
}

func NewTokenSourceManager() *TokenSourceManager {
	return &TokenSourceManager{url: map[Repo]oauth2.TokenSource{}}
}

func (x *TokenSourceManager) SetTokenSource(forRepo Repo, a oauth2.TokenSource) {
	x.lk.Lock()
	defer x.lk.Unlock()
	x.url[forRepo] = a
}

func (x *TokenSourceManager) GetTokenSource(forRepo Repo) oauth2.TokenSource {
	x.lk.Lock()
	defer x.lk.Unlock()
	return x.url[forRepo]
}
