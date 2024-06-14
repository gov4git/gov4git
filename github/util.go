package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/gov4git/gov4git/v2"
	"github.com/gov4git/gov4git/v2/materials"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func FetchRepoMaintainers(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
) []string {

	opts := &github.ListCollaboratorsOptions{}
	users, _, err := ghc.Repositories.ListCollaborators(ctx, repo.Owner, repo.Name, opts)
	must.NoError(ctx, err)

	m := []string{}
	for _, u := range users {
		if u.GetPermissions()["maintainer"] || u.GetPermissions()["admin"] {
			m = append(m, strings.ToLower(u.GetLogin()))
		}
	}
	return m
}

func fetchOpenIssues(ctx context.Context, repo Repo, ghc *github.Client, labelled ...string) []*github.Issue {
	opt := &github.IssueListByRepoOptions{State: "open", Labels: labelled}
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
	repo Repo,
	ghc *github.Client,
	issue *github.Issue,
	subject string,
	payload string,
) {
	replyToIssue(ctx, repo, ghc, issue.GetNumber(), subject, payload)
	closeIssue(ctx, repo, ghc, issue.GetNumber())
}

func replyToIssue(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	issueNum int,
	subject string,
	payload string,
) {

	header := fmt.Sprintf(
		"## <a href=%q><img src=%q alt=\"This project is governed with Gov4Git.\" width=\"65\" /></a> %s\n"+
			"On `%s` by Gov4Git `%s`\n\n",
		materials.Gov4GitWebsiteURL,
		materials.Gov4GitAvatarURL,
		subject,
		time.Now().Format(time.RFC850),
		gov4git.GetVersionInfo().Version,
	)

	comment := &github.IssueComment{
		Body: github.String(header + payload),
	}
	_, _, err := ghc.Issues.CreateComment(ctx, repo.Owner, repo.Name, issueNum, comment)
	must.NoError(ctx, err)
}

const (
	FollowUpSubject = "Follow up"
)

func closeIssue(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	issueNumber int,
) {
	req := &github.IssueRequest{
		State: github.String("closed"),
	}
	_, _, err := ghc.Issues.Edit(ctx, repo.Owner, repo.Name, issueNumber, req)
	must.NoError(ctx, err)
}

func fetchIssueComments(
	ctx context.Context,
	repo Repo,
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

func isJoinApprovalPresent(_ context.Context, approvers []string, comments []*github.IssueComment) bool {
	for _, comment := range comments {
		u := comment.GetUser()
		if u == nil {
			continue
		}
		if !util.IsIn(strings.ToLower(u.GetLogin()), approvers...) {
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

func getIssueAuthorLogin(issue *github.Issue) (string, error) {
	u := issue.GetUser()
	if u == nil {
		return "", fmt.Errorf("github issue without author")
	}
	login := strings.ToLower(u.GetLogin())
	if login == "" {
		return "", fmt.Errorf("github issue author has no login")
	}
	return login, nil
}
