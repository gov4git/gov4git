package tools

import (
	"context"
	"sync"

	"github.com/google/go-github/v58/github"
	govgh "github.com/gov4git/gov4git/v2/github"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
	"golang.org/x/oauth2"
)

func ClearComments(
	ctx context.Context,
	token string,
	repo govgh.Repo,
	issueNo int64,

) {

	// create authenticated GitHub client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	ghc := github.NewClient(tc)

	const maxGoroutines = 5
	throttle := make(chan struct{}, maxGoroutines)

	for {
		opts := &github.IssueListCommentsOptions{}
		comments, _, err := ghc.Issues.ListComments(ctx, repo.Owner, repo.Name, int(issueNo), opts)
		must.NoError(ctx, err)
		if len(comments) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, comment := range comments {
			throttle <- struct{}{}
			wg.Add(1)
			go func(comment *github.IssueComment) {
				defer wg.Done()
				defer func() { <-throttle }()
				base.Infof("Deleting comment %v from issue %v", comment.GetID(), issueNo)
				_, err := ghc.Issues.DeleteComment(ctx, repo.Owner, repo.Name, comment.GetID())
				must.NoError(ctx, err)
			}(comment)
		}
		wg.Wait()
	}
}
