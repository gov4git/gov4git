package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func FetchRepoMaintainers(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
) []string {

	opts := &github.ListCollaboratorsOptions{}
	users, _, err := ghc.Repositories.ListCollaborators(ctx, repo.Owner, repo.Name, opts)
	must.NoError(ctx, err)

	m := []string{}
	for _, u := range users {
		if u.GetPermissions()["maintainer"] || u.GetPermissions()["admin"] {
			m = append(m, u.GetLogin())
		}
	}
	return m
}

func fetchOpenIssues(ctx context.Context, repo GithubRepo, ghc *github.Client, labelled string) []*github.Issue {
	opt := &github.IssueListByRepoOptions{State: "open", Labels: []string{labelled}}
	var allIssues []*github.Issue
	for {
		issues, resp, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
		must.NoError(ctx, err)
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues
}

func replyAndCloseIssue(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
	payload string,
) {
	replyToIssue(ctx, repo, ghc, issue, payload)
	closeIssue(ctx, repo, ghc, issue)
}

func replyToIssue(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
	payload string,
) {

	comment := &github.IssueComment{
		Body: github.String(payload),
	}
	_, _, err := ghc.Issues.CreateComment(ctx, repo.Owner, repo.Name, issue.GetNumber(), comment)
	must.NoError(ctx, err)
}

func closeIssue(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
) {
	req := &github.IssueRequest{
		State: github.String("closed"),
	}
	_, _, err := ghc.Issues.Edit(ctx, repo.Owner, repo.Name, issue.GetNumber(), req)
	must.NoError(ctx, err)
}

func fetchIssueComments(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
) []*github.IssueComment {

	if issue.GetComments() == 0 {
		return nil
	}
	opts := &github.IssueListCommentsOptions{}
	comments, _, err := ghc.Issues.ListComments(ctx, repo.Owner, repo.Name, issue.GetNumber(), opts)
	must.NoError(ctx, err)
	return comments
}

func isJoinApprovalPresent(ctx context.Context, approvers []string, comments []*github.IssueComment) bool {
	for _, comment := range comments {
		u := comment.GetUser()
		if u == nil {
			continue
		}
		if !util.IsIn(u.GetLogin(), approvers...) {
			continue
		}
		// trim empty lines and spaces
		trimmed := strings.ToLower(strings.Trim(comment.GetBody(), ". \t\n\r"))
		if strings.Index(trimmed, JoinRequestApprovalWord) < 0 {
			continue
		}
		return true
	}
	return false
}
